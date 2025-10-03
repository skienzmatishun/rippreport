const puppeteer = require('puppeteer');

async function testCommentSystem() {
    console.log('🚀 Starting live comment system test...');
    
    const browser = await puppeteer.launch({ 
        headless: true,
        args: ['--no-sandbox', '--disable-setuid-sandbox']
    });
    
    try {
        const page = await browser.newPage();
        
        // Enable console logging from the page
        page.on('console', msg => {
            if (msg.type() === 'error') {
                console.log('❌ Browser Error:', msg.text());
            } else if (msg.text().includes('🔄') || msg.text().includes('Comment System')) {
                console.log('📝 Comment System:', msg.text());
            }
        });
        
        // Navigate to the page
        console.log('📄 Loading https://rippreport.com/p/newsletter/');
        await page.goto('https://rippreport.com/p/newsletter/', { 
            waitUntil: 'networkidle2',
            timeout: 30000 
        });
        
        // Wait for comment system to initialize
        console.log('⏳ Waiting for comment system to load...');
        await new Promise(resolve => setTimeout(resolve, 5000));
        
        // Check if comment system container exists
        const commentContainer = await page.$('#ai-comment-system');
        if (!commentContainer) {
            console.log('❌ FAIL: Comment system container not found');
            return;
        }
        console.log('✅ Comment system container found');
        
        // Check if comments are loaded
        const commentsLoaded = await page.evaluate(() => {
            const container = document.querySelector('#ai-comment-system');
            const commentsList = container?.querySelector('#comments-list');
            const comments = commentsList?.querySelectorAll('.comment-thread');
            
            return {
                containerExists: !!container,
                commentsListExists: !!commentsList,
                commentCount: comments ? comments.length : 0,
                containerHTML: container ? container.innerHTML.substring(0, 500) : 'Not found',
                hasLoadingState: container?.innerHTML.includes('Loading comments'),
                hasErrorState: container?.innerHTML.includes('error') || container?.innerHTML.includes('Error'),
                hasComments: comments && comments.length > 0
            };
        });
        
        console.log('📊 Comments Status:', commentsLoaded);
        
        if (commentsLoaded.hasLoadingState) {
            console.log('⏳ Still loading, waiting longer...');
            await new Promise(resolve => setTimeout(resolve, 10000));
            
            // Check again
            const finalStatus = await page.evaluate(() => {
                const container = document.querySelector('#ai-comment-system');
                const commentsList = container?.querySelector('#comments-list');
                const comments = commentsList?.querySelectorAll('.comment-thread');
                const replies = commentsList?.querySelectorAll('.comment-replies .comment');
                
                return {
                    commentCount: comments ? comments.length : 0,
                    replyCount: replies ? replies.length : 0,
                    hasLoadingState: container?.innerHTML.includes('Loading comments'),
                    hasErrorState: container?.innerHTML.includes('error') || container?.innerHTML.includes('Error'),
                    containerHTML: container ? container.innerHTML.substring(0, 1000) : 'Not found'
                };
            });
            
            console.log('📊 Final Status:', finalStatus);
        }
        
        // Test switching to Relevant mode
        console.log('🔄 Testing Relevant mode switch...');
        const relevantButton = await page.$('[data-mode="similarity"]');
        if (relevantButton) {
            console.log('✅ Relevant button found, clicking...');
            await relevantButton.click();
            await new Promise(resolve => setTimeout(resolve, 3000));
            
            const afterSwitch = await page.evaluate(() => {
                const container = document.querySelector('#ai-comment-system');
                const commentsList = container?.querySelector('#comments-list');
                const comments = commentsList?.querySelectorAll('.comment-thread');
                const replies = commentsList?.querySelectorAll('.comment-replies .comment');
                
                return {
                    commentCount: comments ? comments.length : 0,
                    replyCount: replies ? replies.length : 0,
                    isEmpty: !comments || comments.length === 0,
                    containerHTML: container ? container.innerHTML.substring(0, 500) : 'Not found'
                };
            });
            
            console.log('📊 After switching to Relevant:', afterSwitch);
            
            if (afterSwitch.isEmpty) {
                console.log('❌ FAIL: Comments disappeared when switching to Relevant mode');
            } else {
                console.log('✅ PASS: Comments preserved when switching to Relevant mode');
            }
        } else {
            console.log('❌ Relevant button not found');
        }
        
        // Test switching back to Recent
        console.log('🔄 Testing Recent mode switch...');
        const recentButton = await page.$('[data-mode="chronological"]');
        if (recentButton) {
            await recentButton.click();
            await new Promise(resolve => setTimeout(resolve, 3000));
            
            const afterSwitchBack = await page.evaluate(() => {
                const container = document.querySelector('#ai-comment-system');
                const commentsList = container?.querySelector('#comments-list');
                const comments = commentsList?.querySelectorAll('.comment-thread');
                const replies = commentsList?.querySelectorAll('.comment-replies .comment');
                
                return {
                    commentCount: comments ? comments.length : 0,
                    replyCount: replies ? replies.length : 0,
                    isEmpty: !comments || comments.length === 0
                };
            });
            
            console.log('📊 After switching back to Recent:', afterSwitchBack);
        }
        
        // Check for JavaScript errors
        const jsErrors = await page.evaluate(() => {
            return window.commentSystemErrors || [];
        });
        
        if (jsErrors.length > 0) {
            console.log('❌ JavaScript Errors:', jsErrors);
        }
        
        // Final summary
        console.log('\n📋 TEST SUMMARY:');
        console.log('================');
        
        const finalCheck = await page.evaluate(() => {
            const container = document.querySelector('#ai-comment-system');
            const commentsList = container?.querySelector('#comments-list');
            const comments = commentsList?.querySelectorAll('.comment-thread');
            const replies = commentsList?.querySelectorAll('.comment-replies .comment');
            const replyButtons = commentsList?.querySelectorAll('.reply-btn');
            
            return {
                containerExists: !!container,
                commentsVisible: comments && comments.length > 0,
                repliesVisible: replies && replies.length > 0,
                replyButtonsVisible: replyButtons && replyButtons.length > 0,
                totalComments: comments ? comments.length : 0,
                totalReplies: replies ? replies.length : 0,
                hasError: container?.innerHTML.includes('error') || container?.innerHTML.includes('Error')
            };
        });
        
        console.log(`✅ Container exists: ${finalCheck.containerExists}`);
        console.log(`✅ Comments visible: ${finalCheck.commentsVisible} (${finalCheck.totalComments} comments)`);
        console.log(`✅ Replies visible: ${finalCheck.repliesVisible} (${finalCheck.totalReplies} replies)`);
        console.log(`✅ Reply buttons visible: ${finalCheck.replyButtonsVisible}`);
        console.log(`❌ Has errors: ${finalCheck.hasError}`);
        
        if (!finalCheck.commentsVisible) {
            console.log('\n❌ MAJOR ISSUE: No comments are visible on the page');
        } else if (!finalCheck.repliesVisible) {
            console.log('\n❌ MAJOR ISSUE: Comments are visible but replies are not showing');
        } else {
            console.log('\n✅ SUCCESS: Comments and replies are both visible');
        }
        
    } catch (error) {
        console.error('❌ Test failed:', error);
    } finally {
        await browser.close();
    }
}

testCommentSystem().catch(console.error);