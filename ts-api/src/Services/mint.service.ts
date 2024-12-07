import { Mint } from '../Models/mints'

export class mintService {

    async createMint(data: any) {
        try {
            const newMint = await Mint.create(data)
            return newMint

        } catch (error) {
            console.log(error)
        }
    }
}

export const mintServices = new mintService()