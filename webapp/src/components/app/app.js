import  { HomePage, SignInPage, SignUpPage } from '../pages';

import Footer from '../footer/index';
import Header from '../header/index';
import React from 'react';

const App = () => {
    return(
        <div className="flex flex-col items-center font-sans bg-purple-100">
            <Header />
            <HomePage/>
            <Footer />
        </div>
    )
}

export default App;
