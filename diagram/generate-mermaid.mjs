import puppeteer from 'puppeteer';
import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';
import { icons as logos } from '@iconify-json/logos';
import { icons as mdi } from '@iconify-json/mdi';
import { icons as skillIcones } from '@iconify-json/skill-icons';
import { icons as akarIcones } from '@iconify-json/akar-icons';

// Définir __dirname pour les modules ESM
const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// Paths for source and destination directories
const SRC_DIR = path.resolve(__dirname, 'src');
const DST_DIR = path.resolve(__dirname, 'dst');
const CONFIG_PATH = path.resolve(__dirname, 'mermaid-config.json');

// Set default width and height
const DEFAULT_WIDTH = 1920;
const DEFAULT_HEIGHT = 400;

// Combine multiple icon packs
const iconPacks = [
  {
    name: logos.prefix,
    icons: logos,
  },
  {
    name: mdi.prefix,
    icons: mdi,
  },
  {
    name: skillIcones.prefix,
    icons: skillIcones,
  },
  {
    name: akarIcones.prefix,
    icons: akarIcones,
  },
];

// Ensure the destination directory exists
if (!fs.existsSync(DST_DIR)) {
  fs.mkdirSync(DST_DIR, { recursive: true });
}

// Read all .mmd files from the src directory
const mermaidFiles = fs.readdirSync(SRC_DIR).filter((file) => file.endsWith('.mmd'));

// Generate diagrams for all .mmd files
(async () => {
  const browser = await puppeteer.launch();
  const page = await browser.newPage();

  await page.setViewport({
    width: DEFAULT_WIDTH,
    height: DEFAULT_HEIGHT,
  });

  for (const file of mermaidFiles) {
    const filePath = path.join(SRC_DIR, file);
    const fileContent = fs.readFileSync(filePath, 'utf-8');
    const outputFileName = file.replace('.mmd', '.svg');
    const outputPath = path.join(DST_DIR, outputFileName);

    console.log(`Generating diagram for ${file} -> ${outputPath}`);

    // Load a blank page
    await page.goto('data:text/html;charset=UTF-8,<html><head></head><body></body></html>');

    // Load Mermaid.js and configure it
    await page.addScriptTag({
      url: 'https://cdn.jsdelivr.net/npm/mermaid@11.4.0/dist/mermaid.min.js',
    });

    const config = JSON.parse(fs.readFileSync(CONFIG_PATH, 'utf8'));
    await page.evaluate((config, iconPacks) => {
      // Initialize Mermaid with config
      window.mermaid.initialize(config);

      // Register all icon packs
      iconPacks.forEach((pack) => {
        window.mermaid.registerIconPacks([pack]);
      });
    }, config, iconPacks);

    // Render the Mermaid diagram
    const content = await page.evaluate((mermaidCode) => {
      try {
        const element = document.createElement('div');
        element.className = 'mermaid';
        element.innerHTML = mermaidCode;
        document.body.appendChild(element);
        window.mermaid.contentLoaded();
        return new Promise((resolve) => {
          setTimeout(() => {
            resolve(element.innerHTML);
          }, 1000); // Wait for the diagram to render
        });
      } catch (error) {
        console.error('Error rendering Mermaid diagram:', error);
        throw error;
      }
    }, fileContent);

    // Save the rendered SVG content to a file
    fs.writeFileSync(outputPath, content);
  }

  console.log('All diagrams have been generated.');
  await browser.close();
})();
