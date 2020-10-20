const Router = require('./router')
const { LinksTransformer, ProfileTransformer, SocialTransformer, TitleTransformer, BodyTransformer } = require('./transformers')
const config = require("./constants.json")


function getLinks(request) {
    const init = {
        headers: { "content-type": "application/json" },
    }
    const body = JSON.stringify(config.links)
    return new Response(body, init)
}

async function getDefault(request) {
    let pt = new ProfileTransformer(config.avatar, config.name)
    let defaultPage = await fetch(config.pageUrl, {method: "get"})
    let updatedPage = new HTMLRewriter()
        .on("div#links", new LinksTransformer(config.links))
        .on("div#profile", pt)
        .on("img#avatar", pt)
        .on("h1#name", pt)
        .on("div#social", new SocialTransformer(config.socialLinks))
        .on("title", new TitleTransformer(config.name))
        .on("body", new BodyTransformer(config.backgroundColor))
        .transform(defaultPage)

    return updatedPage
}

async function handleRequest(request) {
    const r = new Router()
    r.get('/links', request => getLinks(request))
    r.get('/', request => getDefault(request))
    const resp = await r.route(request)

    return resp
}

addEventListener("fetch", event => {
  event.respondWith(handleRequest(event.request))
})