import { mintServices } from '../Services/mint.service'
import { MintManager } from '../Services/mint.helper'
import { Request, Response } from 'express'
import {MintschemaValidate} from '../Models/mints'

class mintController {
    
    mint = async (req: Request, res: Response) => {
        const data = {
            image: req.body.image,
            title: req.body.title,
            description: req.body.description,
        }
        const {error, value} = MintschemaValidate.validate(data)

        if(error){
            res.send(error.message)

        }else{
            const mintManager = new MintManager()
            const mintpub = await mintManager.mint(data.title, data.description, data.image)

            //call the create post function in the service and pass the data from the request
            value.mint = mintpub
            const mint = await mintServices.createMint(value)
            res.status(201).send(mint)          
        }
    }

}

export const MintController = new mintController()