import express from "express";
import { MintController } from '../Controllers/mint.controller'

//initiating the router
export const router = express.Router()

router.post('/', MintController.mint)