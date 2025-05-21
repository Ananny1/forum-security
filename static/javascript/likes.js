function likePost(postID) {
    console.log("Liking post:", postID);

    fetch('/like', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: new URLSearchParams({ postID: postID }),
    })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                document.getElementById('like-count-' + postID).textContent = data.likeCount;
                document.getElementById('dislike-count-' + postID).textContent = data.dislikeCount;

                const likeButton = document.getElementById('like-' + postID)
                const disLikeButton = document.getElementById('dislike-' + postID)
                if (data.likedByUser) {
                    likeButton.innerHTML = "<i class='bx bxs-like'></i>"
                    disLikeButton.innerHTML = "<i class='bx bx-dislike'></i>"
                } else {
                    likeButton.innerHTML = "<i class='bx bx-like'></i>"
                }
            } else {
                console.log("Post ", postID, " couldn't like.");
            }
        })
        .catch(error => {
            console.error('Error:', error);
        });
}

function dislikePost(postID) {
    console.log("Disliking post:", postID);

    fetch('/dislike', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: new URLSearchParams({ postID: postID }),
    })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                document.getElementById('dislike-count-' + postID).textContent = data.dislikeCount;
                document.getElementById('like-count-' + postID).textContent = data.likeCount;

                const likeButton = document.getElementById('like-' + postID)
                const disLikeButton = document.getElementById('dislike-' + postID)
                if (data.dislikedbyuser) {
                    likeButton.innerHTML = "<i class='bx bx-like'></i>"
                    disLikeButton.innerHTML = "<i class='bx bxs-dislike'></i>"
                } else {
                    disLikeButton.innerHTML = "<i class='bx bx-dislike'></i>"
                }
            } else {
                console.log("Post ", postID, " couldn't dislike.");
            }
        })
        .catch(error => {
            console.error('Error:', error);
        });
}

function likeComment(commentID) {
    console.log("Liking comment:", commentID);

    fetch('/commentlike', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: new URLSearchParams({ commentID: commentID }),
    })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                document.getElementById('cdislike-count-' + commentID).textContent = data.dislikeCount;
                document.getElementById('clike-count-' + commentID).textContent = data.likeCount;

                const likeButton = document.getElementById('clike-' + commentID)
                const disLikeButton = document.getElementById('cdislike-' + commentID)
                if (data.likedByUser) {
                    likeButton.innerHTML = "<i class='bx bxs-like'></i>"
                    disLikeButton.innerHTML = "<i class='bx bx-dislike'></i>"
                } else {
                    likeButton.innerHTML = "<i class='bx bx-like'></i>"
                }
            } else {
                console.log("Comment ", commentID, " couldn't like.");
            }
        })
        .catch(error => {
            console.error('Error:', error);
        });
}

function dislikeComment(commentID) {
    console.log("Disliking comment:", commentID);

    fetch('/commentdislike', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: new URLSearchParams({ commentID: commentID }),
    })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                document.getElementById('cdislike-count-' + commentID).textContent = data.dislikeCount;
                document.getElementById('clike-count-' + commentID).textContent = data.likeCount;

                const likeButton = document.getElementById('clike-' + commentID)
                const disLikeButton = document.getElementById('cdislike-' + commentID)
                if (data.dislikedbyuser) {
                    likeButton.innerHTML = "<i class='bx bx-like'></i>"
                    disLikeButton.innerHTML = "<i class='bx bxs-dislike'></i>"
                } else {
                    disLikeButton.innerHTML = "<i class='bx bx-dislike'></i>"
                }
            } else {
                console.log("Comment ", commentID, " couldn't disliked.");
            }
        })
        .catch(error => {
            console.error('Error:', error);
        });
}
