document.querySelector('.newpost').addEventListener('click', function(event) {
    if (event.target === this || event.target.tagName === 'I') {
        this.classList.toggle('active');
    }
});

document.getElementById('postForm').addEventListener('submit', function(event) {
    const title = document.getElementById('title').value;
    const content = document.getElementById('content').value;
    const cheak = document.querySelectorAll('input[name="category"]');
    let isChecked = Array.from(cheak).some(checkbox => checkbox.checked);


    if (title.length > 50) {
        event.preventDefault();
        alert('Title must be at least 5 characters long and no longer than 50.');
        return; 
    }


    if (content.length > 500) {
        event.preventDefault();
        alert('Content must be at least 20 characters long and no longer than 500.');
        return;
    }

    if (!content.replace(/\s/g, '').length) {
        event.preventDefault();
        alert('content does not accept only spaces');
        return;
    }

    if (!isChecked) {
        event.preventDefault();
        alert('You must select at least one category.');
        return;
    }
});

document.querySelector('.dropdown-button').addEventListener('click', function() {
    const dropdownContent = document.querySelector('.dropdown-content');
    dropdownContent.style.display = dropdownContent.style.display === 'block' ? 'none' : 'block';
});