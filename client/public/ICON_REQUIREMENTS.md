# Icon Requirements for Homescreen App

To prevent pixelation when users save your web app to their device homescreen, you need to create high-resolution icon files.

## Required Icon Files

Place these files in the `/public` directory:

### For Android/PWA (Web App Manifest):
1. **icon-192x192.png** - Minimum size for Android devices
2. **icon-512x512.png** - High-resolution icon for Android and PWA

### For iOS (Apple Touch Icons):
1. **apple-touch-icon-180x180.png** - Required for modern iOS devices (iPhone 6 Plus and newer)
2. **apple-touch-icon-152x152.png** - For iPad (Retina)
3. **apple-touch-icon-120x120.png** - For iPhone (Retina)
4. **apple-touch-icon.png** - Fallback icon (180x180 recommended)

## How to Create Icons

1. Start with a high-resolution source image (at least 1024x1024 pixels)
2. Use an image editor or online tool to create PNG files in the sizes listed above
3. Ensure icons are square (1:1 aspect ratio)
4. Icons should have transparent backgrounds or solid backgrounds that match your brand
5. For best results, design icons specifically for small sizes (avoid fine details)

## Recommended Tools

- **Online**: [PWA Asset Generator](https://github.com/onderceylan/pwa-asset-generator)
- **Online**: [RealFaviconGenerator](https://realfavicongenerator.net/)
- **Image Editor**: Photoshop, GIMP, or any image editor that supports batch resizing

## Testing

After adding the icons:
1. Clear your browser cache
2. On mobile: Add to Home Screen and check the icon quality
3. On desktop: Check the browser tab icon and PWA installation icon

