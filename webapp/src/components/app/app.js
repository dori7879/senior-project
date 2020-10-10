import  { HomePage, SignInPage, SignUpPage } from '../pages';
import {Route, Switch} from 'react-router-dom';

import React from 'react';

const App = () => {
    return(
        <div className="flex flex-col items-center font-sans bg-purple-100">
            <main role='main' >
                <Switch>
                    <Route path='/' exact component={HomePage} />
                    <Route path='/signin' exact component={SignInPage} />
                    <Route path='/signup' exact component={SignUpPage} />
                </Switch>
            </main>
        </div>
    )
}

export default App;
