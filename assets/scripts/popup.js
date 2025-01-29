import { fetchResponse} from "./tools.js"
export { Popup }


const commentSex = document.getElementById("comments-section");
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
            const cmntGrp = document.getElementById('comment-group') 
            cmntGrp.innerHTML = `<textarea placeholder="Type a comment..." type="text" id="comment-field"></textarea>
                        <button class="new-comment" id="${postid}"><i class="fas fa-paper-plane"></i></button>`
            
            const newCmntBtn = document.getElementById(`${postid}`)
            // console.log(newCmntBtn)

            newCmntBtn.addEventListener('click', async()=>{
                const newCmnt = {
                    postid: postid,
                    userid: 1, // to be handled
                    content : document.getElementById('comment-field').value, 
                }
                // console.log('hre=>>>>',newCmnt)
                // console.log(typeof newCmnt.postid)
                await createComment(newCmnt)
                
            })
            
            await displaycomment(postid);
            
        }
      };

      const createComment = async (newCmnt) => {
          console.log(newCmnt)
        const msg= await fetchResponse('/createComment', newCmnt)
        console.log(msg)
        // comment created
        await displaycomment(newCmnt.postId)
      }
   

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
    commentSex.innerHTML=''
    const obj = { ID: postid };
    const cmnts = await fetchResponse(`/comments`, obj);
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