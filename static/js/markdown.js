/**
 * To render markdown content.
 */

document.addEventListener('DOMContentLoaded', function() {
    var markdownContent = document.querySelectorAll('.markdown-body');
    
    markdownContent.forEach(function(element) {
        var rawContent = element.innerHTML;
        var renderedHtml = marked(rawContent);
        element.innerHTML = renderedHtml;
    });
});