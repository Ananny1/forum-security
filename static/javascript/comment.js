function toggleCommentForm(postID) {
    var form = document.getElementById('comment-form-' + postID);
    if (form) {
        if (form.style.display === 'none' || form.style.display === '') {
            form.style.display = 'block';
        } else {
            form.style.display = 'none';
        }
    } else {
        console.error('Comment form with ID ' + 'comment-form-' + postID + ' not found.');
    }
}
document.querySelectorAll('.comment-form').forEach(form => {
    form.addEventListener('submit', function(event) {
        const content = form.querySelector('textarea').value;
        if (!content.replace(/\s/g, '').length) {
            event.preventDefault();
            alert('Comment cannot be empty.');
        }
    });
});

document.querySelectorAll('.comment-form').forEach(form => {
    form.addEventListener('submit', function(event) {
        const content = form.querySelector('textarea').value;
        if ( content.length > 200) {
            event.preventDefault();
            alert('Comment must be at least 5 characters long and no longer than 200.');
        }
    });
});

document.addEventListener('DOMContentLoaded', function() {
    const toggleIcon = document.getElementById('toggleCommentBox');
    const commentForm = document.getElementById('commentForm');

    toggleIcon.addEventListener('click', function() {
        if (commentForm.style.display === 'none' || commentForm.style.display === '') {
            commentForm.style.display = 'block';
        } else {
            commentForm.style.display = 'none';
        }
    });
});
