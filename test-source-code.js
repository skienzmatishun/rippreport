const puppeteer = require('puppeteer');

async function testSourceCode() {
    console.log('🚀 Testing ACTUAL HTML source code on live site');
    
    const browser = await puppeteer.launch({ 
        headless: false, // Show browser so we can see what's happening
        args: ['--no-sandbox', '--disable-setuid-sandbox']
    });
    
    try {
        const page = await browser.newPage();
        
        // Navigate to the LIVE site
        console.log('📄 Loading https://rippreport.com/p/newsletter/');
        await page.goto('https://rippreport.com/p/newsletter/', { 
            waitUntil: 'networkidle2',
            timeout: 30000 
        });
        
        // Wait for comment system to fully load
        await new Promise(resolve => setTimeout(resolve, 8000));
        
        console.log('\n=== CHECKING INITIAL HTML SOURCE ===');
        const initialHTML = await page.evaluate(() => {
            const container = document.querySelector('#ai-comment-system');
            return container ? container.innerHTML : 'NO CONTAINER FOUND';
        });
        
        // Count comments in initial HTML
        const initialCommentCount = (initialHTML.match(/data-comment-id="/g) || []).length;
        const initialReplyCount = (initialHTML.match(/comment-replies/g) || []).length;
        
        console.log(`Initial HTML contains ${initialCommentCount} comments`);
        console.log(`Initial HTML contains ${initialReplyCount} reply containers`);
        console.log('Initial HTML snippet:', initialHTML.substring(0, 500) + '...');
        
        if (initialCommentCount === 0) {
            console.log('❌ NO COMMENTS FOUND IN INITIAL HTML');
            return;
        }
        
        console.log('\n=== CLICKING RELEVANT BUTTON ===');
        const relevantButton = await page.$('[data-mode="similarity"]');
        
        if (!relevantButton) {
            console.log('❌ RELEVANT BUTTON NOT FOUND');
            return;
        }
        
        await relevantButton.click();
        console.log('✅ Clicked Relevant button');
        
        // Wait for the API call and re-rendering
        await new Promise(resolve => setTimeout(resolve, 5000));
        
        console.log('\n=== CHECKING HTML SOURCE AFTER RELEVANT CLICK ===');
        const afterHTML = await page.evaluate(() => {
            const container = document.querySelector('#ai-comment-system');
            return container ? container.innerHTML : 'NO CONTAINER FOUND';
        });
        
        // Count comments in HTML after clicking Relevant
        const afterCommentCount = (afterHTML.match(/data-comment-id="/g) || []).length;
        const afterReplyCount = (afterHTML.match(/comment-replies/g) || []).length;
        
        console.log(`After Relevant HTML contains ${afterCommentCount} comments`);
        console.log(`After Relevant HTML contains ${afterReplyCount} reply containers`);
        console.log('After HTML snippet:', afterHTML.substring(0, 500) + '...');
        
        console.log('\n=== COMPARISON RESULTS ===');
        console.log(`Before: ${initialCommentCount} comments, ${initialReplyCount} reply containers`);
        console.log(`After:  ${afterCommentCount} comments, ${afterReplyCount} reply containers`);
        
        if (afterCommentCount === 0) {
            console.log('❌ CRITICAL ISSUE: ALL COMMENTS DISAPPEARED FROM HTML');
            console.log('❌ The Relevant button is clearing all comments from the DOM');
        } else if (afterCommentCount < initialCommentCount) {
            console.log(`❌ ISSUE: ${initialCommentCount - afterCommentCount} comments disappeared from HTML`);
            console.log('❌ The Relevant button is filtering out some comments');
        } else if (afterCommentCount === initialCommentCount) {
            console.log('✅ SUCCESS: All comments preserved in HTML');
            console.log('✅ The Relevant button is working correctly');
        } else {
            console.log('⚠️ UNEXPECTED: More comments after clicking Relevant?');
        }
        
        // Check if the HTML actually changed
        if (initialHTML === afterHTML) {
            console.log('⚠️ HTML is identical - no changes detected');
        } else {
            console.log('✅ HTML changed - button is functional');
        }
        
        // Keep browser open for manual inspection
        console.log('\n🔍 Browser will stay open for 30 seconds for manual inspection...');
        await new Promise(resolve => setTimeout(resolve, 30000));
        
    } catch (error) {
        console.error('❌ Test failed:', error);
    } finally {
        await browser.close();
    }
}

testSourceCode().catch(console.error);