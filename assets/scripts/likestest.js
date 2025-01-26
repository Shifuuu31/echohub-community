document.addEventListener('DOMContentLoaded', () => {
    console.log('DOM fully loaded and parsed');
  
    initializeLikesDislikes();
  
    const likeButtons = document.querySelectorAll('.like-btn');
    const dislikeButtons = document.querySelectorAll('.dislike-btn');
  
    likeButtons.forEach(button => {
      button.addEventListener('click', () => handleLikeDislike(button, true));
    });
  
    dislikeButtons.forEach(button => {
      button.addEventListener('click', () => handleLikeDislike(button, false));
    });
  });
  
  // Fetch likes/dislikes for a specific (post or comment)
  async function fetchLikesDislikes(entityId, entityType) {
    const response = await fetch(`/likes-dislikes?entityId=${entityId}&entityType=${entityType}`);
    if (!response.ok) {
      throw new Error('Failed to fetch likes/dislikes');
    }
    return response.json();
  }
  
  // Initialize likes/dislikes for all posts and comments
  async function initializeLikesDislikes() {
    const posts = document.querySelectorAll('#post');
    const comments = document.querySelectorAll('.comment');
  
    posts.forEach(async (post) => {
      const likeButton = post.querySelector('.like-btn');
      const entityId = likeButton.getAttribute('data-entity-id');
      const entityType = likeButton.getAttribute('data-entity-type');
  
      const data = await fetchLikesDislikes(entityId, entityType);
  
      const likeCountElement = post.querySelector('.like-count');
      const dislikeCountElement = post.querySelector('.dislike-count');
      if (likeCountElement) likeCountElement.textContent = data.likes;
      if (dislikeCountElement) dislikeCountElement.textContent = data.dislikes;
    });
  
    comments.forEach(async (comment) => {
      const likeButton = comment.querySelector('.like-btn');
      const entityId = likeButton.getAttribute('data-entity-id');
      const entityType = likeButton.getAttribute('data-entity-type');
  
      const data = await fetchLikesDislikes(entityId, entityType);
  
      const likeCountElement = comment.querySelector('.like-count');
      const dislikeCountElement = comment.querySelector('.dislike-count');
      if (likeCountElement) likeCountElement.textContent = data.likes;
      if (dislikeCountElement) dislikeCountElement.textContent = data.dislikes;
    });
  }
  
// like/dislike when click on them
async function handleLikeDislike(button, isLike) {
    const entityId = button.getAttribute('data-entity-id');
    const entityType = button.getAttribute('data-entity-type');
  
    const payload = {
      entityId: parseInt(entityId),
      entityType: entityType,
      liked: isLike,
    };
  
    try {
      const response = await fetch('/like-dislike', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      });
  
      if (!response.ok) throw new Error('Failed to update like/dislike');
      const data = await response.json();
  
      const allEntities = document.querySelectorAll(
        `[data-entity-id="${entityId}"][data-entity-type="${entityType}"]`
      );
  
      allEntities.forEach(entity => {
        const parent = entity.closest('#post') || entity.closest('.comment');
        const likeCountElement = parent.querySelector('.like-count');
        const dislikeCountElement = parent.querySelector('.dislike-count');
  
        if (likeCountElement) likeCountElement.textContent = data.likes;
        if (dislikeCountElement) dislikeCountElement.textContent = data.dislikes;
      });
  
    } catch (error) {
      console.error('Error:', error);
      alert('An error occurred.');
    }
  }