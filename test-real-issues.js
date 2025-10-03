const puppeteer = require('puppeteer');

async function testRealIssues() {
    console.log('🚀 Testing FIXED issues: replies showing, Relevant mode preserving comments');
    
    const browser = await puppeteer.launch({ 
        headless: true, // Run headless for automated testing
        args: ['--no-sandbox', '--disable-setuid-sandbox']
    });
    
    try {
        const page = await browser.newPage();
        
        // Enable console logging
        page.on('console', msg => {
            console.log(`[BROWSER] ${msg.type()}: ${msg.text()}`);
        });
        
        // Navigate to the LIVE site (not development)
        console.log('📄 Loading LIVE SITE: https://rippreport.com/p/newsletter/');
        await page.goto('https://rippreport.com/p/newsletter/', { 
            waitUntil: 'networkidle2',
            timeout: 30000 
        });
        
        // Wait for page to fully load
        await new Promise(resolve => setTimeout(resolve, 8000));
        
        // Check initial state
        console.log('\n=== INITIAL STATE ===');
        const initialState = await page.evaluate(() => {
            const container = document.querySelector('#ai-comment-system');
            const commentsList = container?.querySelector('#comments-list');
            const commentThreads = commentsList?.querySelectorAll('.comment-thread');
            const replies = commentsList?.querySelectorAll('.comment-replies');
            const replyComments = commentsList?.querySelectorAll('.comment-replies .comment');
            
            console.log('Container HTML:', container?.innerHTML.substring(0, 200));
            console.log('Comments list exists:', !!commentsList);
            console.log('Comment threads found:', commentThreads?.length || 0);
            console.log('Reply containers found:', replies?.length || 0);
            console.log('Reply comments found:', replyComments?.length || 0);
            
            return {
                hasContainer: !!container,
                hasCommentsList: !!commentsList,
                threadCount: commentThreads?.length || 0,
                replyContainers: replies?.length || 0,
                replyComments: replyComments?.length || 0,
                isEmpty: !commentThreads || commentThreads.length === 0
            };
        });
        
        console.log('Initial state:', initialState);
        
        if (initialState.threadCount === 0) {
            console.log('❌ ISSUE CONFIRMED: No comments showing at all');
        } else if (initialState.replyComments === 0) {
            console.log('❌ ISSUE CONFIRMED: Comments showing but NO REPLIES visible');
        } else {
            console.log('✅ Comments and replies are showing initially');
        }
        
        // Test clicking Relevant button
        console.log('\n=== TESTING RELEVANT BUTTON CLICK ===');
        const relevantButton = await page.$('[data-mode="similarity"]');
        
        if (relevantButton) {
            console.log('Clicking Relevant button...');
            await relevantButton.click();
            
            // Wait a bit for the transition
            await new Promise(resolve => setTimeout(resolve, 5000));
            
            const afterRelevant = await page.evaluate(() => {
                const container = document.querySelector('#ai-comment-system');
                const commentsList = container?.querySelector('#comments-list');
                const commentThreads = commentsList?.querySelectorAll('.comment-thread');
                const replies = commentsList?.querySelectorAll('.comment-replies .comment');
                
                return {
                    threadCount: commentThreads?.length || 0,
                    replyComments: replies?.length || 0,
                    isEmpty: !commentThreads || commentThreads.length === 0,
                    htmlContent: commentsList?.innerHTML.substring(0, 500) || 'No content'
                };
            });
            
            console.log('After clicking Relevant:', afterRelevant);
            
            if (afterRelevant.isEmpty) {
                console.log('❌ ISSUE CONFIRMED: Relevant mode CLEARED ALL COMMENTS');
            } else if (afterRelevant.replyComments === 0) {
                console.log('❌ ISSUE CONFIRMED: Relevant mode shows comments but NO REPLIES');
            } else {
                console.log('✅ Relevant mode preserved comments and replies');
            }
        } else {
            console.log('❌ Relevant button not found');
        }
        
        // Check for JavaScript errors
        const errors = await page.evaluate(() => {
            return window.errors || [];
        });
        
        if (errors.length > 0) {
            console.log('JavaScript errors:', errors);
        }
        
        // Final verification
        console.log('\n=== FINAL VERIFICATION ===');
        const finalCheck = await page.evaluate(() => {
            const container = document.querySelector('#ai-comment-system');
            const commentsList = container?.querySelector('#comments-list');
            const commentThreads = commentsList?.querySelectorAll('.comment-thread');
            const replies = commentsList?.querySelectorAll('.comment-replies .comment');
            
            return {
                totalThreads: commentThreads?.length || 0,
                totalReplies: replies?.length || 0,
                hasAllComments: (commentThreads?.length || 0) >= 2, // Should have at least 2 root comments
                hasAllReplies: (replies?.length || 0) >= 3 // Should have at least 3 replies
            };
        });
        
        console.log('Final verification:', finalCheck);
        
        if (finalCheck.hasAllComments && finalCheck.hasAllReplies) {
            console.log('✅ SUCCESS: All fixes working correctly!');
            console.log(`✅ Found ${finalCheck.totalThreads} comments and ${finalCheck.totalReplies} replies`);
        } else {
            console.log('❌ FAILURE: Issues still exist');
            console.log(`❌ Expected: ≥2 comments, ≥3 replies. Got: ${finalCheck.totalThreads} comments, ${finalCheck.totalReplies} replies`);
        }
        
    } catch (error) {
        console.error('❌ Test failed:', error);
    } finally {
        await browser.close();
    }
}

testRealIssues().catch(console.error);