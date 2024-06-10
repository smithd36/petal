document.addEventListener('DOMContentLoaded', function() {
    var navbarToggler = document.querySelector('.navbar-toggler');
    var navbarCollapse = document.querySelector('#navbarNav');
    var bsCollapse = new bootstrap.Collapse(navbarCollapse, {
        toggle: false
    });

    navbarToggler.addEventListener('click', function() {
        if (navbarCollapse.classList.contains('show')) {
            bsCollapse.hide();
        } else {
            bsCollapse.show();
        }
    });
});