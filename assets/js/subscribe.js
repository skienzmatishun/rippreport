document.getElementById('subscribe-form').addEventListener('submit', async function(e) {
  e.preventDefault();
  
  const form = this;
  const email = document.getElementById('subscribe-email').value;
  const honeypot = form.querySelector('[name="email_check"]').value;
  const messageEl = document.getElementById('subscribe-message');
  const submitBtn = form.querySelector('button[type="submit"]');
  
  // Honeypot check
  if (honeypot) return;
  
  // Disable button during submission
  submitBtn.disabled = true;
  submitBtn.textContent = 'Subscribing...';
  messageEl.style.display = 'none';
  
  try {
    const response = await fetch('https://802322c7.sibforms.com/serve/MUIFAK0KonyzRnk4yAJ4RwuIvFrwbYaSaKKVyy5aeUCvtAfM-aM7VIBcIhcLMVV6RLPEJgu-MKQO1y1O26svyx98C_maBMo2w_AveJaxlKw5J4xotuf9wljkerGFxTiQ3erBJ2X5DWhSnPNSJ_63lspDcb2XbF51bSg_YM4qxJNqNOWxx3xtijbPgi7Wgwff6-UlByAE3gwuELz8', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded',
      },
      body: new URLSearchParams({
        'EMAIL': email,
        'locale': 'en'
      })
    });
    
    // Brevo returns HTML, so we check if request succeeded
    if (response.ok) {
      messageEl.textContent = 'Thanks for subscribing!';
      messageEl.className = 'subscribe-form__message subscribe-form__message--success';
      messageEl.style.display = 'block';
      form.reset();
    } else {
      throw new Error('Subscription failed');
    }
  } catch (error) {
    messageEl.textContent = 'Something went wrong. Please try again.';
    messageEl.className = 'subscribe-form__message subscribe-form__message--error';
    messageEl.style.display = 'block';
  } finally {
    submitBtn.disabled = false;
    submitBtn.textContent = 'Subscribe';
  }
});