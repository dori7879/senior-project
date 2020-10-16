import  { HomePage, HomeworkPage, SignInPage, SignUpPage, AttendancePage, QuizPage } from '../pages';
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
                    <Route path='/homework' exact component={HomeworkPage} />
                    <Route path='/attendance' exact component={AttendancePage} />
                    <Route path='/quiz' exact component={QuizPage} />
                </Switch>
            </main>
        </div>
    )
}

export default App;
