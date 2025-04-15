#!/bin/bash

# Create necessary directories
mkdir -p static/{css,js,images}

# Copy images
cp -r ../app/assets/images/* static/images/

# Copy CSS files
cp -r ../app/assets/stylesheets/* static/css/

# Copy JavaScript files
cp -r ../app/assets/javascripts/* static/js/

# Create a basic application.css that imports all other CSS files
cat > static/css/application.css << 'EOL'
/* Import all CSS files */
@import 'framework.css';
@import 'utilities.css';
@import 'application_dark.css';
@import 'application_utilities.css';
@import 'application_utilities_dark.css';
@import 'behaviors.css';
@import 'errors.css';
@import 'fonts.css';
@import 'performance_bar.css';
@import 'print.css';
@import 'snippets.css';
@import 'tailwind.css';
@import 'test_environment.css';
@import 'utilities.css';
EOL

# Create a basic application.js that imports all other JS files
cat > static/js/application.js << 'EOL'
// Import all JavaScript files
// Note: In a real application, you would use a bundler like webpack
// This is just a placeholder for demonstration
console.log('GitLab application initialized');
EOL

echo "Assets copied successfully!"
