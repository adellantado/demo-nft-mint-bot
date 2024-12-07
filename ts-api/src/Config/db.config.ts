import dotenv from 'dotenv'
import mongoose from 'mongoose'

dotenv.config()

  //details from the env
const username = process.env.DB_USERNAME
const password = process.env.DB_PASSWORD
const dbName = 'Mints'

//connection string to mongo atlas
// const connectionString = `mongodb+srv://${username}:${password}@cluster0.tjh8e.mongodb.net/${dbName}?retryWrites=true&w=majority`
// const options = {
//     autoIndex: false, // Don't build indexes
//     maxPoolSize: 10, // Maintain up to 10 socket connections
//     serverSelectionTimeoutMS: 5000, // Keep trying to send operations for 5 seconds
//     socketTimeoutMS: 45000, // Close sockets after 45 seconds of inactivity
//     family: 4 // Use IPv4, skip trying IPv6
//   };

// localhost connection
const connectionString = `mongodb://localhost:27017/${dbName}`
const options = {}

export const db = mongoose.connect(connectionString, options)
.then(res => {
    if(res){
        console.log(`Database connection succeffully to ${dbName}`)
    }
    
}).catch(err => {
    console.log(err)
})