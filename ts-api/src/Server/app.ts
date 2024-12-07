import express from 'express'
import { db} from '../Config/db.config'
import { router } from '../Routes/mint.routes'

const app = express()

//middlewares
app.use(express.json())
app.use(express.urlencoded({extended: true}))


//routes
app.use('/api/v1/mint', router)

//db connection then server connection
db.then(() => {
    app.listen(7070, () => console.log('Server is listening on port 7070'))
})
