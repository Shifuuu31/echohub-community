document.addEventListener("DOMContentLoaded", () => {
    DropDown()  
})

const DropDown = () => {
    const profilePictureContainer = document.getElementById('profile-picture-container')
    const dropdown = document.getElementById('dropdown')
    if (profilePictureContainer && dropdown) {
        let isDropdownOpen = false

        const toggleDropdown = (event) => {
            event.stopPropagation()
            isDropdownOpen = !isDropdownOpen
            dropdown.classList.toggle('active')
        }

        const closeDropdown = (event) => {
            if (isDropdownOpen && !dropdown.contains(event.target) && !profilePictureContainer.contains(event.target)) {
                isDropdownOpen = false
                dropdown.classList.remove('active')
            }
        }

        profilePictureContainer.addEventListener('click', toggleDropdown)
        document.addEventListener('click', closeDropdown)

        dropdown.addEventListener('click', (event) => {
            event.stopPropagation()
        })
    }
}