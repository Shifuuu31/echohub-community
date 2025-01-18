export { Popup }

const postContent = document.querySelector("#post .post-info-2")
const commentButton = document.querySelector("#post .post-info-3 button:nth-child(3)")
const popupBackground = document.querySelector(".popup-backgroud")
const closeButton = document.querySelector("#popup-comment button")


const OpenPopup = () => {
    document.querySelector(".popup-backgroud").style.display = "block"
    document.getElementById("popup").style.display = "block"
    document.body.style.overflow = "hidden"
}

const ClosePopup = () => {
    document.querySelector(".popup-backgroud").style.display = "none"
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
