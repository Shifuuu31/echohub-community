export { Popup }

const postContent = document.querySelector("#post .post-body")
const commentButton = document.querySelector("#post .post-categories button:nth-child(3)")
const popupBackground = document.querySelector(".popup-background")
const closeButton = document.querySelector("#user-info-and-buttons button")


const OpenPopup = () => {
    document.querySelector(".popup-background").style.display = "block"
    document.getElementById("popup").style.display = "block"
    document.body.style.overflow = "hidden"
}

const ClosePopup = () => {
    document.querySelector(".popup-background").style.display = "none"
    document.getElementById("popup").style.display = "none"
    document.body.style.overflow = "auto"
}

const Popup = () => {
        postContent.addEventListener("click", OpenPopup)
        commentButton.addEventListener("click", OpenPopup)

        popupBackground.addEventListener("click", (event) => {
            if (event.target === popupBackground) {
                ClosePopup()
            }
        })
        closeButton.addEventListener("click", ClosePopup)
}
