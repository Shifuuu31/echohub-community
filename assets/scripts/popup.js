import { fetchResponse} from "./tools.js"
export { Popup }


const Popup = () => {
    const popup = document.getElementById("popup")
    const popupBackground = document.getElementById("popup-background")
    const closeButton = document.querySelector(".close")

    const attachEventListeners = () => {
        const postsBtns = document.querySelectorAll("#commentBtn, #post-title")
        postsBtns.forEach(postBtn => {
            postBtn.removeEventListener("click", openPopup)
            postBtn.addEventListener("click", (event) => {

                openPopup(event)
            })
        })
    }


    const openPopup = async (event) => {
        let popupPost = document.querySelector("#popup #post");
        if (popup && popupBackground) {
          popupBackground.style.display = popup.style.display = "block";
          const targetedPost = event.target.closest("#post");
          const postid = targetedPost.getAttribute("post-id");
          popupPost.replaceWith(targetedPost.cloneNode(true));
          await displaycomment(postid);
        }
      };
   

    const closePopup = (event) => {
        if (event.target === popupBackground || event.target === closeButton) {
            popupBackground.style.display = popup.style.display = "none"

        }
    }

    if (popupBackground) {
        popupBackground.addEventListener("click", closePopup)
    }

    return attachEventListeners
}


const displaycomment = async (postid) => {
    const commentSex = document.getElementById("comments-section");
    const obj = { ID: postid };
    const cmnts = await fetchResponse(`/comments`, obj);
  commentSex.innerHTML=''
    for (let cmnt of cmnts) {
      let comment = document.createElement("div");
      comment.id = "comment";
      comment.innerHTML = `
                          <div id="user-info-and-buttons">
                              <div id="user-comment-info">
                                  <img src="/assets/imgs/avatar.png" alt="User Avatar" loading="lazy">
                                  <h3>${cmnt.UserName} <br><span>${cmnt.CreationDate}</span></h3>
                              </div>
                          </div>
                          <div id="user-comment-info">
                              <p>${cmnt.CommentContent}</p>
                          </div>`;
      commentSex.appendChild(comment);
    }
}