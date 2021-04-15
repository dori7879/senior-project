import { Redirect, useParams, Link } from 'react-router-dom'
import { useSelector } from 'react-redux'
import { useEffect, useState } from 'react'
import Footer from '../components/Footer'
import Header from '../components/Header'
import QuizService from '../services/quiz';

const ListQuizSubmissions = () => {
  const { randomStr } = useParams()
  const role = useSelector((state) => state.auth.role)
  const isLoggedIn = useSelector((state) => state.auth.isLoggedIn)

  // Quiz related state
  // eslint-disable-next-line no-unused-vars
  const [quizID, setQuizID] = useState(-1) 
  const [quizSubs, setQuizSubs] = useState([])
  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  // eslint-disable-next-line no-unused-vars
  const [closeDate, setCloseDate] = useState(new Date())
  // eslint-disable-next-line no-unused-vars
  const [questions, setQuestions] = useState([])

  // Form related state
  // eslint-disable-next-line no-unused-vars
  const [attachments, setAttachments] = useState([])

  useEffect(() => {
    QuizService.fetchTeacherQuiz(randomStr)
      .then((response) => {
        if (response.data) {
          setQuizID(response.data.ID)
          setTitle(response.data.Title)
          setDescription(response.data.Content)
          setCloseDate(response.data.ClosedAt)
          if ('Submissions' in response.data) {
            setQuizSubs(response.data.Submissions)
          }
          setQuestions(response.data.Questions)
        }
      })
      .catch((error) => {
        alert(error.message)
      })
  }, [randomStr])

  const ckEditorRemoveTags = (data) => {
    const editedData = data.replace('<p>', '').replace('</p>', '')
    return editedData
  }

  if (isLoggedIn && role === 'student') {
    return <Redirect to='/' />
  }

  const data = ckEditorRemoveTags(description)
  const isEmptyDesc = description.trim() === ''
  const isEmptyFiles = attachments.length === 0
  const submittedQuizzes = quizSubs.length > 0 ? quizSubs.length : 0

  return (
    <div>
      <Header />
      <div className='flex flex-col items-center min-h-screen px-4 py-2 bg-purple-100 justify-top sm:px-6 lg:px-8 '>
        <div className='flex flex-col items-center justify-center pb-4'></div>
        <div className='pt-4 text-xl font-bold text-purple-900'>
          Quiz Submissions
        </div>

        <div className='w-3/4 p-4 bg-purple-300 border border-purple-300 rounded'>
          <div className='flex flex-row'>
            <div className='flex flex-col w-3/4'>
              <h2 className='block px-4 pt-1 mb-2 text-xs font-bold tracking-wide text-gray-700 uppercase'>
                <strong>Quiz Title:</strong>{' '}
                <span className='text-purple-900'>{title}</span>
              </h2>

              {isEmptyDesc ? null : (
                <h2 className='block px-4 pt-1 mb-2 text-xs font-bold tracking-wide text-gray-700 uppercase'>
                  <strong>Description:</strong>
                  <br></br> <span className='text-purple-900'>{data}</span>
                </h2>
              )}

              {isEmptyFiles ? null : (
                <h2 className='block px-4 pt-1 mb-2 text-xs font-bold tracking-wide text-gray-700 uppercase'>
                  <strong>Attachments:</strong>
                  <br></br> <span className='text-purple-900'>{attachments}</span>
                </h2>
              )}

              <div className='flex flex-col ml-2 items '>
                {quizSubs.map((sub, index) => (
                  <div
                    key={sub.ID}
                    className='flex flex-col mb-2 border border-purple-700 rounded'
                  >
                    <Link className="block px-4 pt-1 tracking-wide text-gray-900" to={{
                      pathname: "/homeworks/view",
                      state: { id: sub.ID, questions: questions }
                    }}>
                      View Submission
                    </Link>
                    <div className='flex flex-col pb-2 items'>
                      <p className='block px-4 pt-1 text-xs tracking-wide text-gray-700'>
                        <strong>Student Name:</strong>{' '}
                        <span className='text-purple-900'>
                          {sub.StudentFullName ? sub.StudentFullName :
                            sub.Student != null ? sub.Student.FirstName + " " + sub.Student.LastName :
                            ""
                          }
                        </span>
                      </p>
                      <p className='block px-4 pt-1 mb-2 text-xs tracking-wide text-gray-700'>
                        <strong>Submitted at:</strong>{' '}
                        <span className='text-purple-900'>
                          {sub.SubmittedAt}
                        </span>
                      </p>
                      {'Grade' in sub && sub.Grade !== 0 ? (
                        <p className='block px-4 pt-1 text-xs tracking-wide text-gray-700'>
                          <strong>Grade:</strong>
                          <span className='text-purple-900'>
                            {sub.Grade}
                          </span>
                        </p>
                      ) : null}
                      {'Comments' in sub && sub.Comments.trim() === '' ? null : (
                        <p className='block px-4 pt-1 text-xs tracking-wide text-gray-700'>
                          <strong>Comments:</strong>
                          <span className='text-purple-900'>
                            {sub.Comments}
                          </span>
                        </p>
                      )}
                    </div>
                  </div>
                ))}
              </div>

            </div>

            <div className='flex flex-col '>
              <p className='block px-4 pt-1 text-xs tracking-wide text-gray-700'>
                <strong>
                  <span className='font-bold text-purple-800'>
                    {submittedQuizzes}
                  </span>{' '}
                  Quizzes Submitted
                </strong>{' '}
              </p>
              <p className='block px-4 pt-1 text-xs tracking-wide text-gray-700'>
                <strong></strong>{' '}
              </p>
            </div>
          </div>
        </div>
      </div>
      <Footer />
    </div>
  )
}

export default ListQuizSubmissions
