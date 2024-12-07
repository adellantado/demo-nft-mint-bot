import { createNft, mplTokenMetadata } from '@metaplex-foundation/mpl-token-metadata'
import { createUmi } from '@metaplex-foundation/umi-bundle-defaults'
import { createGenericFile, createSignerFromKeypair, generateSigner, keypairIdentity, KeypairSigner, percentAmount, sol, Umi } from '@metaplex-foundation/umi'
import { mockStorage } from '@metaplex-foundation/umi-storage-mock'
import * as fs from 'fs'
import * as path from 'path';
import dotenv from 'dotenv'

dotenv.config({ path: __dirname+'/../../../.env' });

class MintHelper {

    quicknodeRpc: string
    umi: Umi
    creator: KeypairSigner
    nftDetail: {
        name: string,
        symbol: string,
        uri: string,
        royalties: number,
        description: string,
        imgType: string,
        attributes: [
            { trait_type: string, value: string },
        ]
    }

    constructor(quicknodeRpc: string) {
        this.quicknodeRpc = quicknodeRpc
        this.umi= createUmi(this.quicknodeRpc)
        const secret = require(process.env.WALLET_JSON!);
        const creatorWallet = this.umi.eddsa.createKeypairFromSecretKey(new Uint8Array(secret))
        this.creator = createSignerFromKeypair(this.umi, creatorWallet)
        this.umi.use(keypairIdentity(this.creator))
        this.umi.use(mplTokenMetadata())
        this.umi.use(mockStorage())
        this.nftDetail = {
            name: "NFT",
            symbol: "NFT",
            uri: "IPFS_URL_OF_METADATA",
            royalties: 2,
            description: 'NFT description here!',
            imgType: 'image/jpg',
            attributes: [
                {trait_type: 'mint_by', value: 'go_nft_bot'},
            ]
        }
    }

    /**
     * setNftDetails
     */
    public setNftDetails(name: string, description: string, imageName: string) {
        this.nftDetail.name = name
        this.nftDetail.symbol = name.split(' ').map(word => word.charAt(0).toUpperCase()).join('')
        this.nftDetail.description = description
        this.nftDetail.imgType = 'image/'+this.getFileExtensionPath(imageName)
    }

    public async uploadImage(imgName: string): Promise<string> {
        try {
            const imgDirectory = process.env.IMAGE_FOLDER || './uploads'
            const filePath = `${imgDirectory}/${imgName}`
            const fileBuffer = fs.readFileSync(filePath)
            const image = createGenericFile(
                fileBuffer,
                imgName,
                {
                    uniqueName: this.nftDetail.name,
                    contentType: this.nftDetail.imgType
                }
            )
            const [imgUri] = await this.umi.uploader.upload([image])
            this.nftDetail.uri = imgUri
            console.log('Uploaded image:', imgUri)
            return imgUri
        } catch (e) {
            throw e
        }

    }


    public async uploadMetadata(imageUri: string): Promise<string> {
        try {
            const metadata = {
                name: this.nftDetail.name,
                description: this.nftDetail.description,
                image: imageUri,
                attributes: this.nftDetail.attributes,
                properties: {
                    files: [
                        {
                            type: this.nftDetail.imgType,
                            uri: imageUri,
                        },
                    ]
                }
            };
            const metadataUri = await this.umi.uploader.uploadJson(metadata)
            console.log('Uploaded metadata:', metadataUri)
            return metadataUri
        } catch (e) {
            throw e
        }
    }


    public async mintNft(metadataUri: string) {
        try {
            const mint = generateSigner(this.umi)
            await createNft(this.umi, {
                mint,
                name: this.nftDetail.name,
                symbol: this.nftDetail.symbol,
                uri: metadataUri,
                sellerFeeBasisPoints: percentAmount(this.nftDetail.royalties),
                creators: [{ address: this.creator.publicKey, verified: true, share: 100 }],
            }).sendAndConfirm(this.umi)
            console.log(`Created NFT: ${mint.publicKey.toString()}`)
            return mint.publicKey.toString()
        } catch (e) {
            throw e
        }
    }

    private getFileExtensionPath(filename: string): string {
        return path.extname(filename).slice(1)
    }

}

class MintManager {

    helper: MintHelper

    constructor() {
        this.helper = new MintHelper(process.env.QUICKNODE_RPC!)
    }

    async mint(name: string, description: string, imageName: string) {
        this.helper.setNftDetails(name, description, imageName)
        const imageUri = await this.helper.uploadImage(imageName)
        const metadataUri = await this.helper.uploadMetadata(imageUri)
        return await this.helper.mintNft(metadataUri)
    }

}

export { MintManager }