import  {Schema, model,} from 'mongoose'
import Joi from 'joi'

export const MintschemaValidate = Joi.object({
    image: Joi.string().required(),
    description: Joi.string().required(),
    title: Joi.string().required(),
})

interface IMints {
    title: string,
    description: string,
    image: string,
    mint: string,
}

const mintSchema = new Schema<IMints>({
    title: {
        type: String,
        required: true
    },

    description: {
        type: String,
        required: true
    },

    image: {
        type: String,
        required: true
    },

    mint: {
        type: String,
        required: false
    },

})

 export const Mint = model<IMints>('Mint', mintSchema )