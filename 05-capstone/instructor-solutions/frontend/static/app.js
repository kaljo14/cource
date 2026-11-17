const form = document.getElementById('form');
const input = document.getElementById('text');
const submit = document.getElementById('submit');
const submitLabel = document.getElementById('submit-label');
const status = document.getElementById('status');
const boardStatus = document.getElementById('board-status');
const list = document.getElementById('notes');
const refresh = document.getElementById('refresh');
let loading = false;

function formMessage(message, error = false) {
  status.textContent = message;
  status.classList.toggle('error', error);
}

function renderNote(note) {
  const card = document.createElement('li');
  card.className = 'note';
  const heading = document.createElement('div');
  heading.className = 'note-heading';
  const avatar = document.createElement('span');
  avatar.className = 'avatar';
  avatar.textContent = '✧';
  avatar.setAttribute('aria-hidden', 'true');
  const author = document.createElement('span');
  author.className = 'note-author';
  author.textContent = 'A fellow learner';
  heading.append(avatar, author);
  const body = document.createElement('p');
  body.className = 'note-body';
  // Never interpret a submitted note as HTML.
  body.textContent = note.text;
  const footer = document.createElement('div');
  footer.className = 'note-footer';
  const timestamp = document.createElement('time');
  const date = new Date(note.createdAt);
  if (!Number.isNaN(date.getTime())) {
    timestamp.dateTime = date.toISOString();
    timestamp.textContent = date.toLocaleString(undefined, { month: 'short', day: 'numeric', hour: 'numeric', minute: '2-digit' });
    timestamp.title = date.toLocaleString();
  } else {
    timestamp.textContent = 'Shared with the class';
  }
  const marker = document.createElement('span');
  marker.className = 'note-marker';
  marker.textContent = '↗';
  marker.setAttribute('aria-hidden', 'true');
  footer.append(timestamp, marker);
  card.append(heading, body, footer);
  return card;
}

async function load() {
  if (loading) return false;
  loading = true;
  refresh.disabled = true;
  list.setAttribute('aria-busy', 'true');
  boardStatus.classList.remove('error');
  boardStatus.textContent = 'Loading the latest notes…';
  try {
    const response = await fetch('/api/notes', { cache: 'no-store' });
    if (!response.ok) throw new Error('The board is temporarily unavailable. Please try Refresh in a moment.');
    const notes = await response.json();
    if (!Array.isArray(notes)) throw new Error('We couldn’t load the notes. Please try again.');
    list.replaceChildren(...notes.map(renderNote));
    document.getElementById('note-count').textContent = notes.length;
    document.getElementById('empty').hidden = notes.length !== 0;
    boardStatus.textContent = notes.length ? 'Latest notes first · up to 100 shared thoughts' : '';
    return true;
  } catch (error) {
    boardStatus.textContent = error.message === 'Failed to fetch' ? 'Connection lost. Your existing notes are still shown; try Refresh when you’re back online.' : error.message;
    boardStatus.classList.add('error');
    return false;
  } finally {
    loading = false;
    refresh.disabled = false;
    list.setAttribute('aria-busy', 'false');
  }
}

async function theme() {
  try {
    const response = await fetch('/config', { cache: 'no-store' });
    if (!response.ok) return;
    const config = await response.json();
    if (!/^#[0-9a-fA-F]{6}$/.test(config.backgroundColor)) return;
    document.body.style.backgroundColor = config.backgroundColor;
    document.getElementById('color-swatch').style.backgroundColor = config.backgroundColor;
    document.getElementById('color').textContent = config.backgroundColor;
  } catch {
    // Keep the last known color during temporary network interruptions.
  }
}

input.addEventListener('input', () => {
  document.getElementById('counter').textContent = `${input.value.length} / 200`;
});
refresh.addEventListener('click', load);
form.addEventListener('submit', async event => {
  event.preventDefault();
  if (submit.disabled) return;
  const text = input.value.trim();
  if (!text) { formMessage('Write a little something before posting.', true); input.focus(); return; }
  submit.disabled = true;
  input.disabled = true;
  submitLabel.textContent = 'Posting your note…';
  formMessage('');
  try {
    const response = await fetch('/api/notes', {
      method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ text })
    });
    if (!response.ok) throw new Error(response.status === 400 ? 'Please keep your note between 1 and 200 characters.' : 'Your note wasn’t saved. Please try again in a moment.');
    input.value = '';
    document.getElementById('counter').textContent = '0 / 200';
    const updated = await load();
    formMessage(updated ? 'Your note is on the board. Thanks for sharing!' : 'Your note was saved. Refresh the board to see it.');
  } catch (error) {
    formMessage(error.message === 'Failed to fetch' ? 'We couldn’t confirm your post. Refresh the board before trying again.' : error.message, true);
  } finally {
    submit.disabled = false;
    input.disabled = false;
    submitLabel.textContent = 'Post to the board';
    input.focus();
  }
});

load();
theme();
setInterval(theme, 5000);
