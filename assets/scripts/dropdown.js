document.addEventListener("DOMContentLoaded", () => {
    DropDown()
})

const DropDown = () => {
    const profilePictureContainer = document.getElementById('profile-picture-container')
    const dropdown = document.getElementById('dropdown')

    const toggleDropdown = (event) => {
        event.stopPropagation()
        dropdown.classList.toggle('active')
    }

    profilePictureContainer.addEventListener('click', toggleDropdown)
    document.addEventListener('click', ()=>{dropdown.classList.remove('active')})
}