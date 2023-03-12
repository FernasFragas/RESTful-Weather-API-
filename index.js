const express = require("express")
const router = express.Router()
const app = express()
const fetch = require("node-fetch")

const PORT = 3000

app.set("view engine", "hbs")
app.use(express.static("public"))
app.use(express.urlencoded( {extended: true} ))
app.use(router)
app.listen(PORT, () => {
    console.log("Listening in " + PORT)
})

router.get("/", (req,res, next) => {
    return fetch('http://localhost:8080/')
        .then(r => r.json())
        .then(resp => {
            res.render("index", {resp: resp})
        })
        .catch(error => console.error(error))
})