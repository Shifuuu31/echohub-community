export { initializeLikeDislike }

//this logic generated using ai for test purposes 
// we'll not use it it's a trash code

function handleLikeDislike(likeButton, dislikeButton) {
    likeButton.addEventListener('click', () => {
        if (likeButton.classList.contains('active')) {
            likeButton.classList.remove('active')
            likeButton.querySelector('img').style.filter = 'invert(50%)'
        } else {
            likeButton.classList.add('active')
            likeButton.querySelector('img').style.filter = 'invert(58%) sepia(61%) saturate(2878%) hue-rotate(230deg) brightness(102%) contrast(101%)'
            likeButton.style.transform = 'scale(1.2)'
            setTimeout(() => {
                likeButton.style.transform = 'scale(1)'
            }, 200)
            dislikeButton.classList.remove('active')
            dislikeButton.querySelector('img').style.filter = 'invert(50%)'
        }
        syncPopupState(likeButton, dislikeButton)
    })

    dislikeButton.addEventListener('click', () => {
        if (dislikeButton.classList.contains('active')) {
            dislikeButton.classList.remove('active')
            dislikeButton.querySelector('img').style.filter = 'invert(50%)'
        } else {
            dislikeButton.classList.add('active')
            dislikeButton.querySelector('img').style.filter = 'invert(58%) sepia(61%) saturate(2878%) hue-rotate(230deg) brightness(102%) contrast(101%)'
            dislikeButton.style.transform = 'scale(1.2)'
            setTimeout(() => {
                dislikeButton.style.transform = 'scale(1)'
            }, 200)
            likeButton.classList.remove('active')
            likeButton.querySelector('img').style.filter = 'invert(50%)'
        }
        syncPopupState(likeButton, dislikeButton)
    })
}

function syncPopupState(likeButton, dislikeButton) {
    const popupLikeButton = document.querySelector('#popup .buttons button:first-child')
    const popupDislikeButton = document.querySelector('#popup .buttons button:nth-child(2)')

    if (likeButton.classList.contains('active')) {
        popupLikeButton.classList.add('active')
        popupLikeButton.querySelector('img').style.filter = 'invert(58%) sepia(61%) saturate(2878%) hue-rotate(230deg) brightness(102%) contrast(101%)'
        popupDislikeButton.classList.remove('active')
        popupDislikeButton.querySelector('img').style.filter = 'invert(50%)'
    } else if (dislikeButton.classList.contains('active')) {
        popupDislikeButton.classList.add('active')
        popupDislikeButton.querySelector('img').style.filter = 'invert(58%) sepia(61%) saturate(2878%) hue-rotate(230deg) brightness(102%) contrast(101%)'
        popupLikeButton.classList.remove('active')
        popupLikeButton.querySelector('img').style.filter = 'invert(50%)'
    } else {
        popupLikeButton.classList.remove('active')
        popupLikeButton.querySelector('img').style.filter = 'invert(50%)'
        popupDislikeButton.classList.remove('active')
        popupDislikeButton.querySelector('img').style.filter = 'invert(50%)'
    }
}

function initializeLikeDislike() {
    const likeButtons = document.querySelectorAll('.buttons button:first-child')
    const dislikeButtons = document.querySelectorAll('.buttons button:nth-child(2)')

    likeButtons.forEach((likeButton, index) => {
        const dislikeButton = dislikeButtons[index]
        handleLikeDislike(likeButton, dislikeButton)
    })

    const popupLikeButtons = document.querySelectorAll('#popup .buttons button:first-child')
    const popupDislikeButtons = document.querySelectorAll('#popup .buttons button:nth-child(2)')

    popupLikeButtons.forEach((popupLikeButton, index) => {
        const popupDislikeButton = popupDislikeButtons[index]
        handleLikeDislike(popupLikeButton, popupDislikeButton)
    })
}
