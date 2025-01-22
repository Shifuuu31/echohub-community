document.addEventListener('DOMContentLoaded', () => {
  console.log('DOM fully loaded and parsed');

  const likeButtons = document.querySelectorAll('.like-btn');
  const dislikeButtons = document.querySelectorAll('.dislike-btn');

  console.log(`Found ${likeButtons.length} like buttons and ${dislikeButtons.length} dislike buttons`);

  likeButtons.forEach(button => {
    button.addEventListener('click', () => handleLikeDislike(button, true));
  });

  dislikeButtons.forEach(button => {
    button.addEventListener('click', () => handleLikeDislike(button, false));
  });
});

function handleLikeDislike(button, isLike) {
  console.log('Button clicked:', button);

  const entityId = button.getAttribute('data-entity-id');
  const entityType = button.getAttribute('data-entity-type');

  const postContainer = button.closest('.post');
  if (!postContainer) {
    console.error('Post container not found');
    return;
  }

  const likeCountElement = postContainer.querySelector('.like-count');
  const dislikeCountElement = postContainer.querySelector('.dislike-count');

  console.log(`Sending request: entityId=${entityId}, entityType=${entityType}, isLike=${isLike}`);
  console.log('Like count element:', likeCountElement);
  console.log('Dislike count element:', dislikeCountElement);

  const payload = {
    entityId: parseInt(entityId), 
    entityType: entityType,
    liked: isLike, 
  };

  console.log('Payload:', payload);

  fetch('/like-dislike', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(payload),
  })
  .then(response => {
    if (!response.ok) {
      return response.text().then(text => {
        throw new Error(`HTTP error! Status: ${response.status}, Response: ${text}`);
      });
    }
    return response.json();
  })
  .then(data => {
    console.log('Response:', data);
    if (data.likes !== undefined && data.dislikes !== undefined) {
      if (likeCountElement) {
        likeCountElement.textContent = data.likes;
      } else {
        console.error('Like count element not found');
      }
      if (dislikeCountElement) {
        dislikeCountElement.textContent = data.dislikes;
      } else {
        console.error('Dislike count element not found');
      }
    } else {
      alert('Invalid response from server');
    }
  })
  .catch(error => {
    console.error('Error:', error);
    alert('An error occurred. Check the console for details.');
  });
}