import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next';
import { useForm } from "react-hook-form";
import Footer from '../components/Footer'
import Header from '../components/Header'
import QuizService from '../services/quiz';

const ViewQuizSubmission = (props) => {
  const id = props.location.state.id
  const questions = props.location.state.questions
  const { t } = useTranslation(['translation', 'quiz']);
  // eslint-disable-next-line no-unused-vars
  const { register, handleSubmit, control, getValues, errors, setValue } = useForm();
  
  // Quiz related state
  const [grade, setGrade] = useState('')
  const [comments, setComments] = useState('')
  const [studentFullName, setStudentFullName] = useState('')
  const [quizID, setQuizID] = useState(-1) 
  // eslint-disable-next-line no-unused-vars
  const [submittedDate, setSubmittedDate] = useState(new Date())
  const [responses, setResponses] = useState([])

  // Form related state
  const [successful, setSuccessful] = useState(false)

  useEffect(() => {
    QuizService.fetchQuizSubmission(id)
      .then((response) => {
        if (response.data) {
          setQuizID(response.data.QuizID)
          setGrade(response.data.Grade)
          setComments(response.data.Comments)
          setSubmittedDate(response.data.SubmittedAt)
          if (response.data.StudentFullName !== "") {
            setStudentFullName(response.data.StudentFullName)
          } else if (response.data.Student != null) {
            setStudentFullName(response.data.User.FirstName + " " + response.data.User.LastName)
          }
          setResponses(response.data.Responses)
        }
      })
      .catch((error) => {
        alert(error.message)
      })
  }, [id, register]);

  const onSubmit = (data) =>  {
    QuizService.updateQuizSubmission({ ...data, SubID: id, QuizID: quizID, Grade: parseFloat(data.Grade) })
      .then(
        (response) => {
          setSuccessful(true);
        },
        (error) => {
          alert(error.message);
        }
      )
  };
  
  return (
    <div>
      <Header />
      <div className='flex flex-col items-center min-h-screen px-4 py-2 bg-purple-100 justify-top sm:px-6 lg:px-8 '>
        <div className='flex flex-col items-center justify-center w-3/4 pb-4'>
          <div className='pt-2 text-xl font-bold text-purple-900'>{t('quiz:quizsubmission', 'Quiz Submission')}</div>
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
                  <span className='text-purple-900'>{JSON.stringify(submittedDate)}</span>
                </h2>
              </div>
            </div>
          </div>
          <div className='w-full p-4 bg-purple-300 border border-purple-300 rounded-b'>
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
                    Overall Comments*:
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

                {questions.map((question, index) => {
                  switch (question.Type) {
                    case 'single':
                      return (
                        <li key={question.ID}>
                          <div className='flex flex-col'>
                            <div className='flex flex-row'>
                              <div className='flex flex-col w-full'>
                                <div className='flex flex-row justify-between w-full px-4 pt-1 mb-2 text-xs font-bold tracking-wide text-gray-700'>
                                  <div className='flex flex-row'>
                                    <label className='block px-2 pt-1 mb-2 text-xs font-bold tracking-wide text-purple-800'>
                                      Question
                                    </label>
                                    <div className='px-2 py-1 text-xs leading-tight text-gray-700 border border-purple-400 rounded w-80 focus:outline-none focus:bg-white'>
                                      {question.Content}
                                    </div>
                                  </div>
                                </div>
                                <div className='pb-2 ml-4 text-xs text-gray-700'>
                                  Tick the correct answer
                                </div>
                                {question.Choices.map((choice, idx) => (
                                  <div key={idx} className='flex flex-row'>
                                    <div className='mx-2 outline-none'>
                                      {responses[index].QuestionID === question.ID ?
                                          idx === responses[index].SingleChoiceResponse &&
                                          responses[index].SingleChoiceResponse === question.SingleChoiceAnswer ?
                                          'Correct' :
                                          'Wrong'
                                        : null
                                      }
                                    </div>
                                    <div className='px-2 py-1 text-xs leading-tight text-gray-700 focus:outline-none focus:bg-white'>
                                      {choice}
                                    </div>
                                    <input
                                      name={`Responses[${index}].Grade`}
                                      ref={register}
                                      className='px-1 mb-3 mr-4 text-xs leading-tight text-gray-700 border border-purple-400 rounded focus:outline-none focus:bg-white'
                                      type='text'
                                      defaultValue={responses[index].Grade}
                                      placeholder={t('homework:entergrade','Enter grade')}
                                    />
                                    <input
                                      name={`Responses[${index}].Comments`}
                                      ref={register}
                                      className='px-1 mb-3 mr-4 text-xs leading-tight text-gray-700 border border-purple-400 rounded focus:outline-none focus:bg-white'
                                      type='text'
                                      defaultValue={responses[index].Comments}
                                      placeholder={t('homework:entercomments','Enter comments')}
                                    />
                                  </div>
                                ))}
                              </div>
                            </div>
                            <input
                              ref={register()}
                              control={control}
                              type='hidden'
                              name={`Responses[${index}].QuestionID`}
                              className='outline-none'
                              defaultValue={question.ID}
                            />
                          </div>
                        </li>
                      );
                    case 'multiple':
                      return (
                        <li key={question.ID}>
                          <div className='flex flex-col'>
                            <div className='flex flex-row'>
                              <div className='flex flex-col w-full'>
                                <div className='flex flex-row justify-between w-full px-4 pt-1 mb-2 text-xs font-bold tracking-wide text-gray-700'>
                                  <div className='flex flex-row'>
                                    <label className='block px-2 pt-1 mb-2 text-xs font-bold tracking-wide text-purple-800'>
                                      Question
                                    </label>
                                    <div className='px-2 py-1 text-xs leading-tight text-gray-700 border border-purple-400 rounded w-80 focus:outline-none focus:bg-white'>
                                      {question.Content}
                                    </div>
                                  </div>
                                </div>
                                <div className='pb-2 ml-4 text-xs text-gray-700'>
                                  Tick correct answers
                                </div>
                                {question.Choices.map((choice, idx) => (
                                  <div key={idx} className='flex flex-row'>
                                    <div className='mx-2 outline-none'>
                                      {responses[index].QuestionID === question.ID ?
                                          responses[index].MultipleChoiceResponse.includes(idx) &&
                                          question.MultipleChoiceAnswer.includes(idx) ?
                                          'Correct' :
                                          'Wrong'
                                        : null
                                      }
                                    </div>
                                    <div className='px-2 py-1 text-xs leading-tight text-gray-700 focus:outline-none focus:bg-white'>
                                      {choice}
                                    </div>
                                    <input
                                      name={`Responses[${index}].Grade`}
                                      ref={register}
                                      className='px-1 mb-3 mr-4 text-xs leading-tight text-gray-700 border border-purple-400 rounded focus:outline-none focus:bg-white'
                                      type='text'
                                      defaultValue={responses[index].Grade}
                                      placeholder={t('homework:entergrade','Enter grade')}
                                    />
                                    <input
                                      name={`Responses[${index}].Comments`}
                                      ref={register}
                                      className='px-1 mb-3 mr-4 text-xs leading-tight text-gray-700 border border-purple-400 rounded focus:outline-none focus:bg-white'
                                      type='text'
                                      defaultValue={responses[index].Comments}
                                      placeholder={t('homework:entercomments','Enter comments')}
                                    />
                                  </div>
                                ))}
                              </div>
                            </div>
                            <input
                              ref={register()}
                              control={control}
                              type='hidden'
                              name={`Responses[${index}].QuestionID`}
                              className='outline-none'
                              defaultValue={question.ID}
                            />
                          </div>
                        </li>
                      );
                    case 'truefalse':
                      return (
                        <li key={question.ID}>
                          <div className='flex flex-col w-full'>
                            <div className='flex flex-row justify-between px-4 pt-1 mb-2 text-xs font-bold tracking-wide text-gray-700'>
                              <div className='flex flex-row'>
                                <label className='block px-2 pt-1 mb-2 text-xs font-bold tracking-wide text-purple-800'>
                                  Question
                                </label>
                                <div className='px-2 py-1 text-xs leading-tight text-gray-700 border border-purple-400 rounded w-80 focus:outline-none focus:bg-white'>
                                  {question.Content}
                                </div>
                              </div>
                            </div>
                            <div className='ml-4'>
                              <div className='flex flex-row'>
                                <div className='mx-2 outline-none'>
                                  {responses[index].QuestionID === question.ID ?
                                      responses[index].TrueFalseResponse === questions[index].TrueFalseAnswer ?
                                      'Correct' :
                                      'Wrong'
                                    : null
                                  }
                                </div>
                              
                                <div>
                                  <input
                                    type='radio'
                                    name={`Responses[${index}].TrueFalseResponse`}
                                    checked={responses[index].TrueFalseResponse === true}
                                    className='outline-none'
                                    disabled
                                  />
                                  <label
                                    htmlFor='choice'
                                    className='px-1 pt-1 mb-2 text-xs font-bold tracking-wide text-gray-700 uppercase '
                                  >
                                    true
                                  </label>
                                  <br></br>
                                  <input
                                    type='radio'
                                    name={`Responses[${index}].TrueFalseResponse`}
                                    checked={responses[index].TrueFalseResponse === false}
                                    className='outline-none'
                                    disabled
                                  />
                                  <label
                                    htmlFor='choice'
                                    className='px-1 pt-1 mb-2 text-xs font-bold tracking-wide text-gray-700 uppercase '
                                  >
                                    false
                                  </label>
                                  <br></br>
                                </div>
                              </div>
                            </div>
                            <input
                              name={`Responses[${index}].Grade`}
                              ref={register}
                              className='px-1 mb-3 mr-4 text-xs leading-tight text-gray-700 border border-purple-400 rounded focus:outline-none focus:bg-white'
                              type='text'
                              defaultValue={responses[index].Grade}
                              placeholder={t('homework:entergrade','Enter grade')}
                            />
                            <input
                              name={`Responses[${index}].Comments`}
                              ref={register}
                              className='px-1 mb-3 mr-4 text-xs leading-tight text-gray-700 border border-purple-400 rounded focus:outline-none focus:bg-white'
                              type='text'
                              defaultValue={responses[index].Comments}
                              placeholder={t('homework:entercomments','Enter comments')}
                            />
                            <input
                              ref={register()}
                              control={control}
                              type='hidden'
                              name={`Responses[${index}].QuestionID`}
                              className='outline-none'
                              defaultValue={question.ID}
                            />
                          </div>
                        </li>
                      );
                    case 'open':
                      return (
                        <li key={question.ID}>
                          <div className='flex flex-col w-full'>
                            <div className='flex flex-row justify-between px-4 pt-1 mb-2 text-xs font-bold tracking-wide text-gray-700'>
                              <div className='flex flex-row'>
                                <label className='block px-2 pt-1 mb-2 text-xs font-bold tracking-wide text-purple-800'>
                                  Question
                                </label>
                                <div className='px-2 py-1 text-xs leading-tight text-gray-700 border border-purple-400 rounded w-80 focus:outline-none focus:bg-white'>
                                  {question.Content}
                                </div>
                              </div>
                            </div>
                            {questions[index].OpenAnswer ? 
                            <div className='flex flex-row'>
                              <label className='block px-4 pt-1 mb-2 text-xs font-bold tracking-wide text-gray-700'>
                                Teacher's Answer
                              </label>
                              <div className='px-2 py-1 text-xs leading-tight text-gray-700 border border-purple-400 rounded focus:outline-none focus:bg-white'>
                                {questions[index].OpenAnswer}
                              </div>
                            </div> 
                            : null
                            }
                            <div className='flex flex-row'>
                              <label className='block px-4 pt-1 mb-2 text-xs font-bold tracking-wide text-gray-700'>
                                Student Response
                              </label>
                              <div className='px-2 py-1 text-xs leading-tight text-gray-700 border border-purple-400 rounded focus:outline-none focus:bg-white'>
                                {responses[index].OpenResponse}
                              </div>
                            </div>
                            <input
                              name={`Responses[${index}].Grade`}
                              ref={register}
                              className='px-1 mb-3 mr-4 text-xs leading-tight text-gray-700 border border-purple-400 rounded focus:outline-none focus:bg-white'
                              type='text'
                              defaultValue={responses[index].Grade}
                              placeholder={t('homework:entergrade','Enter grade')}
                            />
                            <input
                              name={`Responses[${index}].Comments`}
                              ref={register}
                              className='px-1 mb-3 mr-4 text-xs leading-tight text-gray-700 border border-purple-400 rounded focus:outline-none focus:bg-white'
                              type='text'
                              defaultValue={responses[index].Comments}
                              placeholder={t('homework:entercomments','Enter comments')}
                            />
                            <input
                              ref={register()}
                              control={control}
                              type='hidden'
                              name={`Responses[${index}].QuestionID`}
                              className='outline-none'
                              defaultValue={question.ID}
                            />
                          </div>
                        </li>
                      );
                    default:
                      return null
                }})}

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

export default ViewQuizSubmission