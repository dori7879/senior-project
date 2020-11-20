import App from './components/app/index';
import { Provider } from 'react-redux';
import React from 'react';
import ReactDOM from 'react-dom';
import {BrowserRouter as Router} from 'react-router-dom';
import store from './store';

ReactDOM.render(
        <Provider store={store}>
          <Router>
              <App />
          </Router>,
       </Provider>,
  document.getElementById('root')
)