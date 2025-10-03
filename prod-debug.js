
    const puppeteer = require('puppeteer');

    (async () => {
      const browser = await puppeteer.launch({ headless: true });
      const page = await browser.newPage();
      
      console.log('Setting up event listeners for logs and network requests...');
      
      page.on('console', msg => {
        console.log(`BROWSER CONSOLE: ${msg.text()}`);
      });
      
      page.on('response', async (response) => {
        const request = response.request();
        if (request.url().includes('/api/comments/similar/')) {
          console.log(`--- Network Response ---`);
          console.log(`URL: ${response.url()}`);
          console.log(`Status: ${response.status()}`);
          try {
            const text = await response.text();
            console.log(`Body: ${text.substring(0, 500)}`);
          } catch (e) {
            console.log('Could not read response body.');
          }
          console.log(`----------------------`);
        }
      });

      try {
        console.log('Navigating to https://rippreport.com/p/newsletter/ ...');
        await page.goto('https://rippreport.com/p/newsletter/', { waitUntil: 'networkidle2' });
        
        console.log('Waiting for the comment system to load...');
        await page.waitForSelector('#comments-list .comment-thread', { timeout: 10000 });
        console.log('Comments loaded initially.');

        const initialHtml = await page.$eval('#comments-list', el => el.innerHTML);
        console.log('Initial comments HTML is present.');

        console.log('Clicking the "Relevant" button...');
        await page.click('.ordering-button[data-mode="similarity"]');
        
        console.log('Waiting for 2 seconds for changes to apply...');
        await new Promise(resolve => setTimeout(resolve, 2000));

        const finalHtml = await page.$eval('#comments-list', el => el.innerHTML);
        
        console.log('--- Final State Analysis ---');
        if (finalHtml.trim() === '' || finalHtml.includes('comments-empty')) {
          console.log('RESULT: The comments container is empty. The comments have disappeared.');
        } else {
          console.log('RESULT: Comments are still visible in the container.');
        }
        console.log('--------------------------');

      } catch (error) {
        console.error('An error occurred during the Puppeteer test:', error);
      } finally {
        await browser.close();
      }
    })();
    