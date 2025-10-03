
    const puppeteer = require('puppeteer');

    (async () => {
      const browser = await puppeteer.launch({ headless: true });
      const page = await browser.newPage();
      
      console.log('Setting up listeners...');
      
      page.on('console', msg => {
        console.log(`BROWSER CONSOLE: ${msg.text()}`);
      });
      
      page.on('response', async (response) => {
        const request = response.request();
        if (request.url().includes('/api/comments/')) {
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
        
        console.log('Waiting for 5 seconds to observe initial load...');
        await new Promise(resolve => setTimeout(resolve, 5000));

        const initialHtml = await page.$eval('#comments-list', el => el.innerHTML).catch(() => 'COULD NOT FIND #comments-list');
        
        console.log('--- Initial State Analysis ---');
        if (initialHtml.includes('comment-thread')) {
          console.log('RESULT: Comments loaded correctly on initial page load.');
        } else if (initialHtml.includes('comments-empty')) {
            console.log('RESULT: Comments container is present but shows the "empty" message.');
        } else if (initialHtml.includes('comments-error')) {
            console.log('RESULT: Comments container shows an error message.');
        } else {
          console.log('RESULT: The comments container is empty or missing.');
          console.log('Initial HTML:', initialHtml.substring(0, 500));
        }
        console.log('--------------------------');

      } catch (error) {
        console.error('An error occurred during the Puppeteer test:', error);
      } finally {
        await browser.close();
      }
    })();
    