import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next';
import { useForm } from "react-hook-form";
import moment from "moment";
import Parser from 'html-react-parser'
import Footer from '../components/Footer'
import Header from '../components/Header'
import HomeworkService from '../services/homework';

const ViewHWSubmission = (props) => {
  const id = props.location.state.id
  const { t } = useTranslation(['translation', 'homework']);
  // eslint-disable-next-line no-unused-vars
  const { register, handleSubmit, setValue } = useForm();
  
  // Homework related state
  const [response, setResponse] = useState('')
  const [grade, setGrade] = useState('')
  const [comments, setComments] = useState('')
  const [studentFullName, setStudentFullName] = useState('')
  const [homeworkID, setHomeworkID] = useState(-1)
  // eslint-disable-next-line no-unused-vars
  const [submittedDate, setSubmittedDate] = useState(new Date())

  // Form related state
  const [successful, setSuccessful] = useState(false)

  useEffect(() => {
    HomeworkService.fetchHWSubmission(id)
      .then((response) => {
        if (response.data) {
          setHomeworkID(response.data.HomeworkID)
          setResponse(response.data.Response)
          setGrade(response.data.Grade)
          setComments(response.data.Comments)
          setSubmittedDate(response.data.SubmittedAt)
          if (response.data.StudentFullName !== "") {
            setStudentFullName(response.data.StudentFullName)
          } else if (response.data.Student != null) {
            setStudentFullName(response.data.User.FirstName + " " + response.data.User.LastName)
          }
        }
      })
      .catch((error) => {
        alert(error.message)
      })
  }, [id, register]);

  const onSubmit = (data) =>  {
    HomeworkService.updateHWSubmission({ ...data, SubID: id, HomeworkID: homeworkID, Grade: parseFloat(data.Grade) })
      .then(
        (response) => {
          setSuccessful(true);
        },
        (error) => {
          alert(error.message);
        }
      )
  };
  
  const isEmptyDesc = response.trim() === ''
  
  return (
    <div>
      <Header />
      <div className='flex flex-col items-center min-h-screen px-4 py-2 bg-purple-100 justify-top sm:px-6 lg:px-8 '>
        <div className='flex flex-col items-center justify-center w-3/4 pb-4'>
          <div className='pt-2 text-xl font-bold text-purple-900'>{t('homework:hwsubmission', 'Homework Submission')}</div>
          <div className='flex flex-row items-center w-full border border-purple-300 rounded-t items '>
            <div>
              {studentFullName ? (
                <div>
                  <h2 className='block px-4 pt-1 mb-2 text-xs font-bold tracking-wide text-gray-700 uppercase'>
                    <strong>{t('homework:fullname', 'Full Name*:')}</strong>
                    <br></br>{' '}
                  </h2>
                  <p className='ml-4 text-purple-900'>{studentFullName}</p>
                </div>
              ) : null}
              <div className='flex flex-row items-center items'>
                <h2 className='block px-4 pt-1 mb-2 text-xs font-bold tracking-wide text-gray-700 uppercase'>
                  <strong>{t('homework:submittedat', 'Submitted At:')}</strong>{' '}
                  <span className='text-purple-900'>{moment(submittedDate, moment.ISO_8601).format('lll')}</span>
                </h2>
              </div>
            </div>
          </div>
          <div className='w-full p-4 bg-purple-300 border border-purple-300 rounded-b'>
            {isEmptyDesc ? null : (
              <div>
                <h2 className='block px-4 pt-1 mb-2 text-xs font-bold tracking-wide text-gray-700 uppercase'>
                  <strong>{t('homework:description', 'Description')}</strong>
                  <br></br>{' '}
                </h2>
                <p className='ml-4 text-purple-900'>{Parser(response)}</p>
              </div>
            )}

            {successful ? (
              <div className='form-group'>
                <div className='w-full pt-2 text-xl font-bold text-center text-purple-900 border border-purple-300 rounded'>
                {t('homework:updated', 'Updated')}
                </div>
              </div>
            ) : (
              <form onSubmit={handleSubmit(onSubmit)}>
                <div className='flex flex-row items'>
                  <label className='block px-4 pt-1 mb-4 text-xs font-bold tracking-wide text-gray-700 uppercase'>
                    Grade*:
                  </label>
                  <input
                    name="Grade"
                    ref={register}
                    className='px-1 mb-3 mr-4 text-xs leading-tight text-gray-700 border border-purple-400 rounded focus:outline-none focus:bg-white'
                    type='text'
                    defaultValue={grade}
                    placeholder={t('homework:entergrade','Enter grade')}
                  />
                </div>
                <div className='flex flex-row items'>
                  <label className='block px-4 pt-1 mb-4 text-xs font-bold tracking-wide text-gray-700 uppercase'>
                    Comments*:
                  </label>
                  <input
                    name="Comments"
                    ref={register}
                    className='px-1 mb-3 mr-4 text-xs leading-tight text-gray-700 border border-purple-400 rounded focus:outline-none focus:bg-white'
                    type='text'
                    defaultValue={comments}
                    placeholder={t('homework:entercomments','Enter comments')}
                  />
                </div>
                <div className='flex justify-center'>
                  <button
                    type='submit'
                    className='relative flex justify-center w-full px-2 py-1 mb-2 text-sm font-medium leading-4 text-purple-200 transition duration-150 ease-in-out bg-purple-800 border border-transparent rounded-md hover:bg-purple-500 focus:outline-none'
                  >
                    {t('homework:updatehw', 'Update')}
                  </button>
                </div>
              </form>
            )}
          </div>
        </div>
        <small className=''>{t('homework:required', '* - Required field')} </small>
      </div>
      <Footer />
    </div>
  )
}

export default ViewHWSubmission