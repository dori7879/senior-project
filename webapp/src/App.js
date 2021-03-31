import React from 'react';
import logo from './logo.svg';
import './App.css';
import { useTranslation } from 'react-i18next';

function App() {
  const { t, i18n } = useTranslation(['translation', 'welcome']);

  const changeLanguage = code => {
    i18n.changeLanguage(code);
  };

  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <p className="text-green-600">
          Edit <code>src/App.js</code> and save to reload.
        </p>
        <a
          className="App-link"
          href="https://reactjs.org"
          target="_blank"
          rel="noopener noreferrer"
        >
          {t('welcome:title', 'Learn React')}
        </a>
        <button className="mt-2" type="button" onClick={() => changeLanguage('ru')}>{t('translation:ru')}</button>
        <button type="button" onClick={() => changeLanguage('en')}>{t('translation:en')}</button>
      </header>
    </div>
  );
}

export default App;
