// Go Fork Package Registry - JavaScript

document.addEventListener('DOMContentLoaded', function () {
    // Initialize app
    initializeCopyButtons();
    initializeSearch();
    initializePackageCards();
    initia// Initialize theme management
    function initializeTheme() {
        const prefersDarkScheme = window.matchMedia('(prefers-color-scheme: dark)');
        const themeToggle = document.getElementById('theme-toggle');
        const savedTheme = localStorage.getItem('theme');

        // Apply saved theme if exists
        if (savedTheme) {
            document.body.classList.toggle('dark-theme', savedTheme === 'dark');
            if (themeToggle) themeToggle.checked = (savedTheme === 'dark');
        } else if (prefersDarkScheme.matches) {
            // Use OS preference if no saved theme
            document.body.classList.add('dark-theme');
            if (themeToggle) themeToggle.checked = true;
        }

        // Handle theme toggle click
        if (themeToggle) {
            themeToggle.addEventListener('change', function () {
                document.body.classList.toggle('dark-theme', this.checked);
                localStorage.setItem('theme', this.checked ? 'dark' : 'light');
            });
        }

        // Handle OS theme change
        prefersDarkScheme.addEventListener('change', (e) => {
            // Only apply if user hasn't set a preference
            if (!localStorage.getItem('theme')) {
                document.body.classList.toggle('dark-theme', e.matches);
                if (themeToggle) themeToggle.checked = e.matches;
            }
        });

        // Initialize back to top button
        initializeBackToTop();
    }

    // Back to top functionality
    function initializeBackToTop() {
        const backToTopButton = document.getElementById('back-to-top');

        if (!backToTopButton) return;

        // Show button when scrolling down
        window.addEventListener('scroll', () => {
            if (window.pageYOffset > 300) {
                backToTopButton.style.display = 'flex';
            } else {
                backToTopButton.style.display = 'none';
            }
        });

        // Scroll to top when clicked
        backToTopButton.addEventListener('click', () => {
            window.scrollTo({
                top: 0,
                behavior: 'smooth'
            });
        });
    } initializeTheme();
    initializePerformanceMonitoring();
});

// Copy to clipboard functionality
function initializeCopyButtons() {
    const copyButtons = document.querySelectorAll('.copy-btn');

    copyButtons.forEach(button => {
        button.addEventListener('click', async function () {
            const textToCopy = this.dataset.clipboard;

            try {
                await navigator.clipboard.writeText(textToCopy);
                showCopySuccess(this);
            } catch (err) {
                // Fallback for older browsers
                fallbackCopyToClipboard(textToCopy, this);
            }
        });
    });
}

function showCopySuccess(button) {
    const originalText = button.textContent;

    button.textContent = '✅ Copied!';
    button.classList.add('copied');

    setTimeout(() => {
        button.textContent = originalText;
        button.classList.remove('copied');
    }, 2000);
}

function fallbackCopyToClipboard(text, button) {
    const textArea = document.createElement('textarea');
    textArea.value = text;
    document.body.appendChild(textArea);
    textArea.select();

    try {
        document.execCommand('copy');
        showCopySuccess(button);
    } catch (err) {
        console.error('Copy failed:', err);
        showCopyError(button);
    }

    document.body.removeChild(textArea);
}

function showCopyError(button) {
    const originalText = button.textContent;

    button.textContent = '❌ Failed';
    button.style.background = '#dc3545';

    setTimeout(() => {
        button.textContent = originalText;
        button.style.background = '';
    }, 2000);
}

// Search functionality
function initializeSearch() {
    const searchInput = document.querySelector('.search-input');
    if (!searchInput) return;

    // Add search suggestions
    searchInput.addEventListener('input', debounce(handleSearchInput, 300));

    // Add keyboard shortcuts
    document.addEventListener('keydown', function (e) {
        // Ctrl/Cmd + K to focus search
        if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
            e.preventDefault();
            searchInput.focus();
        }

        // Escape to clear search
        if (e.key === 'Escape' && document.activeElement === searchInput) {
            searchInput.value = '';
            searchInput.blur();
        }
    });
}

function handleSearchInput(e) {
    const query = e.target.value.toLowerCase();
    const packageCards = document.querySelectorAll('.package-card');

    packageCards.forEach(card => {
        const packagePath = card.querySelector('.package-path').textContent.toLowerCase();
        const packageDesc = card.querySelector('.package-description').textContent.toLowerCase();
        const packageTags = Array.from(card.querySelectorAll('.tag')).map(tag => tag.textContent.toLowerCase());

        const matches = packagePath.includes(query) ||
            packageDesc.includes(query) ||
            packageTags.some(tag => tag.includes(query));

        if (matches || query === '') {
            card.style.display = 'block';
            card.style.opacity = '1';
        } else {
            card.style.display = 'none';
            card.style.opacity = '0';
        }
    });

    updatePackageCount();
}

// Package card animations
function initializePackageCards() {
    const packageCards = document.querySelectorAll('.package-card');

    // Add intersection observer for scroll animations
    if (window.IntersectionObserver) {
        const observer = new IntersectionObserver((entries) => {
            entries.forEach(entry => {
                if (entry.isIntersecting) {
                    entry.target.style.opacity = '1';
                    entry.target.style.transform = 'translateY(0)';
                }
            });
        }, {
            threshold: 0.1,
            rootMargin: '0px 0px -50px 0px'
        });

        packageCards.forEach((card, index) => {
            card.style.opacity = '0';
            card.style.transform = 'translateY(20px)';
            card.style.transition = `opacity 0.6s ease ${index * 0.1}s, transform 0.6s ease ${index * 0.1}s`;
            observer.observe(card);
        });
    }

    // Add click tracking
    packageCards.forEach(card => {
        card.addEventListener('click', function (e) {
            if (e.target.tagName !== 'A' && e.target.tagName !== 'BUTTON') {
                const detailLink = this.querySelector('.btn-primary');
                if (detailLink) {
                    window.location.href = detailLink.href;
                }
            }
        });
    });
}

// Stats functionality
function initializeStats() {
    const statNumbers = document.querySelectorAll('.stat-number');

    if (statNumbers.length === 0) return;

    // Animate numbers on load
    statNumbers.forEach(stat => {
        const finalValue = parseInt(stat.textContent);
        animateNumber(stat, 0, finalValue, 1000);
    });
}

function animateNumber(element, start, end, duration) {
    const range = end - start;
    const startTime = performance.now();

    function update(currentTime) {
        const elapsed = currentTime - startTime;
        const progress = Math.min(elapsed / duration, 1);

        const easeOut = 1 - Math.pow(1 - progress, 3);
        const current = Math.floor(start + range * easeOut);

        element.textContent = current;

        if (progress < 1) {
            requestAnimationFrame(update);
        }
    }

    requestAnimationFrame(update);
}

// Update package count for search results
function updatePackageCount() {
    const visibleCards = document.querySelectorAll('.package-card[style*="display: block"], .package-card:not([style*="display: none"])');
    const statNumber = document.querySelector('.stat-number');

    if (statNumber) {
        statNumber.textContent = visibleCards.length;
    }
}

// Utility function for debouncing
function debounce(func, wait) {
    let timeout;
    return function executedFunction(...args) {
        const later = () => {
            clearTimeout(timeout);
            func.apply(this, args);
        };
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
    };
}

// Theme management
function initializeTheme() {
    const prefersDarkScheme = window.matchMedia('(prefers-color-scheme: dark)');

    if (prefersDarkScheme.matches) {
        document.body.classList.add('dark-theme');
    }

    prefersDarkScheme.addEventListener('change', (e) => {
        if (e.matches) {
            document.body.classList.add('dark-theme');
        } else {
            document.body.classList.remove('dark-theme');
        }
    });
}

// Package filtering by category
function filterByCategory(category) {
    const packageCards = document.querySelectorAll('.package-card');

    packageCards.forEach(card => {
        const cardCategory = card.dataset.category;

        if (category === '' || cardCategory === category) {
            card.style.display = 'block';
        } else {
            card.style.display = 'none';
        }
    });

    updatePackageCount();
}

// Enhanced search with highlighting
function highlightSearchTerms(query) {
    const packageCards = document.querySelectorAll('.package-card');

    packageCards.forEach(card => {
        const elements = card.querySelectorAll('.package-path, .package-description, .tag');

        elements.forEach(element => {
            const text = element.textContent;
            const highlightedText = text.replace(
                new RegExp(`(${query})`, 'gi'),
                '<mark>$1</mark>'
            );

            if (text !== highlightedText) {
                element.innerHTML = highlightedText;
            }
        });
    });
}

// Performance monitoring
function initializePerformanceMonitoring() {
    if ('performance' in window) {
        window.addEventListener('load', () => {
            setTimeout(() => {
                const perfData = performance.timing;
                const loadTime = perfData.loadEventEnd - perfData.navigationStart;

                console.log(`Page load time: ${loadTime}ms`);

                // Send to analytics if needed
                if (loadTime > 3000) {
                    console.warn('Page load time is slow:', loadTime);
                }
            }, 0);
        });
    }
}

// Export functions for potential external use
window.GoForkRegistry = {
    filterByCategory,
    highlightSearchTerms,
    updatePackageCount
};
