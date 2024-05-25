/* @refresh reload */
import '@abraham/reflection';
import { render } from 'solid-js/web';

import './di';
import './index.css';
import App from './App';
import { container } from 'tsyringe';
import { BackendApi } from './services/backend-api';

console.log(container);
container.resolve(BackendApi);

const root = document.getElementById('root');

if (import.meta.env.DEV && !(root instanceof HTMLElement)) {
    throw new Error(
        'Root element not found. Did you forget to add it to your index.html? Or maybe the id attribute got misspelled?',
    );
}

render(() => <App />, root!);
