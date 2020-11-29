import  { AttendancePage, HomePage, HomeworkPage, QuizPage, SignInPage, SignUpPage } from '../pages';
import {Route, Switch} from 'react-router-dom';

import LinkPage from '../pages/link-page';
import MyProfilePage from '../pages/profile-page';
import React from 'react';
import StudentHWPage from '../pages/student-hw-page';
import TeacherHWPage from '../pages/teacher-hw-page';

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
                    <Route path='/profile' exact component={MyProfilePage} />
                    <Route path='/link' exact component={LinkPage} />
                    <Route path='/student-hw-page/:randomStr' exact component={StudentHWPage} />
                    <Route path='/student-hw-page' exact component={StudentHWPage} />
                    <Route path='/teacher-hw-page' exact component={TeacherHWPage} />
                    <Route path='/teacher-hw-page/:randomStr' component={TeacherHWPage} />
                </Switch>
            </main>
        </div>
    )
}

export default App;
