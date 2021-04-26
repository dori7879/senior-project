import { Route, Switch } from 'react-router-dom'
import Home from './pages/Home'
import Login from './pages/Login'
import Registration from './pages/Registration'
import CreateHomework from './pages/CreateHomework'
import Profile from './pages/Profile'
import Links from './pages/Links'
import SubmitHomework from './pages/SubmitHomework'
import ViewGroup from './pages/ViewGroup'
import AcceptGroupShare from './pages/AcceptGroupShare'
import ListHWSubmissions from './pages/ListHWSubmissions'
import LinksQuiz from './pages/LinksQuiz'
import CreateQuiz from './pages/CreateQuiz'
import ListQuizSubmissions from './pages/ListQuizSubmissions'
import SubmitQuiz from './pages/SubmitQuiz'
import ViewHWSubmission from './pages/ViewHWSubmission'
import ViewQuizSubmission from './pages/ViewQuizSubmission'
import CreateAttendance from './pages/CreateAttendance'
import SubmitAttendance from './pages/SubmitAttendance'
import ListAttSubmissions from './pages/ListAttSubmissions'
import LinksAttendance from './pages/LinksAttendance'

function App() {
  return (
    <div className='flex flex-col items-center font-sans bg-purple-100'>
      <main role='main'>
        <Switch>
          <Route path='/' exact component={Home} />
          <Route path='/login' exact component={Login} />
          <Route path='/signup' exact component={Registration} />
          {/* Profile and Groups */}
          <Route path='/profile' exact component={Profile} />
          <Route
            path='/groups/:groupID'
            component={ViewGroup}
          />
          <Route
            path='/accept-group-link/:shareHash'
            component={AcceptGroupShare}
          />
          {/* Homework */}
          <Route path='/homeworks' exact component={CreateHomework} />
          <Route
            path='/homeworks/submit/:randomStr'
            component={SubmitHomework}
          />
          <Route path='/homeworks/list-submissions/:randomStr' component={ListHWSubmissions} />
          <Route path='/link' exact render={(props) => <Links {...props}/>} />
          <Route path='/homeworks/view' exact render={(props) => <ViewHWSubmission {...props}/>} />
          {/* Quiz */}
          <Route path='/quizzes' exact component={CreateQuiz} />
          <Route
            path='/quizzes/submit/:randomStr'
            component={SubmitQuiz}
          />
          <Route
            path='/quizzes/list-submissions/:randomStr'
            exact
            component={ListQuizSubmissions}
          />
          <Route path='/link-quiz' exact render={(props) => <LinksQuiz {...props}/>} />
          <Route path='/quizzes/view' exact render={(props) => <ViewQuizSubmission {...props}/>} />
          {/* Attendance */}
          <Route path='/attendances' exact component={CreateAttendance} />
          <Route
            path='/attendances/submit/:randomStr'
            component={SubmitAttendance}
          />
          <Route path='/attendances/list-submissions/:randomStr' component={ListAttSubmissions} />
          <Route path='/link-attendance' exact render={(props) => <LinksAttendance {...props}/>} />
        </Switch>
      </main>
    </div>
  )
}

export default App