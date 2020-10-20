class LinksTransformer {
    constructor(links) { this.links = links }

    async element(element) {
        this.links.forEach((link) => {
            element.append(`<a href="${link.url}">${link.name}</a>`, {html: true})
        })
    }
}

class ProfileTransformer {
	constructor(avatar, name) {
		this.avatar = avatar
		this.name = name
	}

	async element(element) {
		switch(element.tagName) {
			case "div":
				element.removeAttribute("style")
				break;
			case "img":
				element.setAttribute("src", this.avatar)
			case "h1":
				element.setInnerContent(this.name)
		}
	}
}

class SocialTransformer {
	constructor(links) { this.links = links }

	async element(element) {
		element.removeAttribute("style")
		this.links.forEach((link) => {
			element.append(`
				<a href="${link.url}">
					<img src="${link.icon}" />
				</a>
			`, {html: true})
		})
	}
}

class TitleTransformer {
	constructor(name) { this.name = name }
	async element(element) {
		element.setInnerContent(this.name)
	}
}

class BodyTransformer {
	constructor(color) { this.color = color }
	async element(element) {
		element.setAttribute("class", this.color)
	}
}

module.exports = {
	LinksTransformer,
	ProfileTransformer,
	SocialTransformer,
	TitleTransformer,
	BodyTransformer
}