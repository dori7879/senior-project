import  { HomePage, SignInPage, SignUpPage } from '../pages';
import {Route, Switch} from 'react-router-dom';

import Footer from '../footer/index';
import Header from '../header/index';
import React from 'react';

const App = () => {
    return(
        <div className="flex flex-col items-center font-sans bg-purple-100">
            <Header />
            <main role='main' >
                <Switch>
                    <Route path='/' exact component={HomePage} />
                    <Route path='/signin' exact component={SignInPage} />
                    <Route path='/signup' exact component={SignUpPage} />
                </Switch>
            </main>
            <Footer />
        </div>
    )
}

export default App;
