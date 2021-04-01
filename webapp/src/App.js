import { Route, Switch } from 'react-router-dom'
import Home from './pages/Home'
import Login from './pages/Login'
import Registration from './pages/Registration'
import CreateHomework from './pages/CreateHomework'
import Profile from './pages/Profile'
import Links from './pages/Links'
import SubmitHomework from './pages/SubmitHomework'
/*
import Attendance from './pages/Attendance'
import LinksQuiz from './pages/LinksQuiz.js'
import Quiz from './pages/Quiz'
import ViewHomeworks from './pages/ViewHomeworks'
import ViewQuizzes from './pages/ViewQuizzes'*/

function App() {
  return (
    <div className='flex flex-col items-center font-sans bg-purple-100'>
      <main role='main'>
        <Switch>
          <Route path='/' exact component={Home} />
          <Route path='/signin' exact component={Login} />
          <Route path='/signup' exact component={Registration} />
          <Route path='/homework' exact component={CreateHomework} />
          <Route path='/link' exact component={Links} />
          <Route path='/profile' exact component={Profile} />
          <Route path='/student-hw-page' exact component={SubmitHomework} />
          <Route
            path='/student-hw-page/:randomStr'
            component={SubmitHomework}
          />
           {/*<Route path='/attendance' exact component={Attendance} />
          <Route path='/teacher-hw-page/:randomStr' component={ViewHomeworks} />
          <Route path='/quiz' exact component={Quiz} />
          <Route path='/teacher-hw-page' exact component={ViewHomeworks} />
          <Route path='/link_quiz' exact component={LinksQuiz} />
          <Route path='/teacher-quiz-page' exact component={ViewQuizzes} />
          <Route
            path='/teacher-quiz-page/c4gtnV23yui'
            exact
            component={ViewQuizzes}
          />*/}
        </Switch>
      </main>
    </div>
  )
}

export default App