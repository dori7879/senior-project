import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next';
import { useForm, Controller } from "react-hook-form";
import moment from "moment";
import Footer from '../components/Footer'
import Header from '../components/Header'
import QuizService from '../services/quiz';

const ViewQuizSubmission = (props) => {
  const id = props.location.state.id
  const { t } = useTranslation(['translation', 'quiz']);
  // eslint-disable-next-line no-unused-vars
  const { register, watch, handleSubmit, control, getValues, errors, setValue } = useForm();
  
  // Quiz related state
  const [grade, setGrade] = useState(0)
  const [maxGrade, setMaxGrade] = useState(0)
  const [comments, setComments] = useState('')
  const [studentFullName, setStudentFullName] = useState('')
  const [quizID, setQuizID] = useState(-1) 
  // eslint-disable-next-line no-unused-vars
  const [submittedDate, setSubmittedDate] = useState(new Date())
  const [responses, setResponses] = useState([])
  const [questions, setQuestions] = useState([])

  const watchResponses = watch("Responses", [])

  // Form related state
  const [successful, setSuccessful] = useState(false)

  useEffect(() => {
    QuizService.fetchQuizSubmission(id)
      .then((response) => {
        if (response.data) {
          setQuizID(response.data.QuizID)
          setMaxGrade(response.data.Quiz.MaxGrade)
          setGrade(parseFloat(response.data.Grade))
          setComments(response.data.Comments)
          setSubmittedDate(response.data.SubmittedAt)
          if (response.data.StudentFullName !== "") {
            setStudentFullName(response.data.StudentFullName)
          } else if (response.data.Student != null) {
            setStudentFullName(response.data.User.FirstName + " " + response.data.User.LastName)
          }
          setResponses(response.data.Responses)
          setQuestions([...props.location.state.questions])
        }
      })
      .catch((error) => {
        alert(error.message)
      })
  }, [id, register, props.location.state.questions]);

  useEffect(() => {
    // setGrade(watchResponses.reduce((acc, r) => acc + parseInt(r.Grade) < maxGrade ? acc + parseInt(r.Grade) : acc, 0))
    if (watchResponses.length > 0) {
      setGrade(watchResponses.reduce((acc, r) => acc + parseInt(r.Grade), 0))
    }
  }, [watchResponses])

  const onSubmit = (data) =>  {
    let body = {
      ...data,
      SubID: id, 
      QuizID: quizID, 
      Grade: grade,
    }

    QuizService.updateQuizSubmission(body)
      .then(
        (response) => {
          setSuccessful(true);
        },
        (error) => {
          alert(error.message);
        }
      )
  };

  const findResponseDefaultValue = (qID) => {
    let curResponse = responses.find(r => r.QuestionID === qID)
    return responses.length > 0 ? curResponse.ID : 0
  }

  const renderSingleCorrect = (q, r, cIdx) => {
    if (q.ID !== r.QuestionID) {
      return null
    }

    // If student's answer is the current index
    if (r.SingleChoiceResponse === cIdx) {
      // If student selected the correct answer
      if (q.SingleChoiceAnswer === cIdx) {
        return <span>&#9989;</span>
      // If student selected the wrong answer
      } else {
        return <span>&#10060;</span>
      }
    // If the correct answer is the current index and student did not select it
    } else if (q.SingleChoiceAnswer === cIdx) {
      return <span className="text-xs text-green-600">Answer</span>
    }
  }

  const renderMultipleCorrect = (q, r, cIdx) => {
    if (q.ID !== r.QuestionID) {
      return null
    }

    // If student's one of the answers is the current index
    if (r.MultipleChoiceResponse.includes(cIdx)) {
      // If student selected a correct answer
      if (q.MultipleChoiceAnswer.includes(cIdx)) {
        return <span>&#9989;</span>
      // If student selected the wrong answer
      } else {
        return <span>&#10060;</span>
      }
    // If the correct answer is the current index and student did not select it
    } else if (q.MultipleChoiceAnswer.includes(cIdx)) {
      return <span className="text-xs text-green-600">Answer</span>
    }
  }
  
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
                  <span className='text-purple-900'>{moment(submittedDate, moment.ISO_8601).format('lll')}</span>
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
                <div className='flex flex-row items-center'>
                  <label className='block px-4 pt-1 mb-4 text-xs font-bold tracking-wide text-gray-700 uppercase'>
                    Grade*:
                  </label>
                  <div className="px-1 mb-3 mr-4 text-xs leading-tight text-gray-700 border border-purple-400 rounded ">
                    {grade}
                  </div>
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

                <ul>
                {questions.map((question, index) => {
                  switch (question.Type) {
                    case 1:
                      return (
                        <li key={question.ID} className="mt-5">
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
                                {question.Choices.map((choice, idx) => (
                                  <div key={idx} className='flex flex-row items-center px-4'>
                                    <div className='w-12 mx-2 outline-none'>
                                      {renderSingleCorrect(questions[index], responses[index], idx)}
                                    </div>
                                    <div className='px-2 py-1 text-xs leading-tight text-gray-700 focus:outline-none focus:bg-white'>
                                      {choice}
                                    </div>
                                  </div>
                                ))}
                                <div className='flex flex-row items-center px-4 mt-2'>
                                  <span className="ml-2 text-xs font-bold tracking-wide text-gray-700">Grade: </span>
                                  <Controller
                                    control={control}
                                    name={`Responses[${index}].Grade`}
                                    defaultValue={responses[index].Grade}
                                    render={(props) => (
                                      <input
                                        type='text'
                                        name={`Responses[${index}].Grade`}
                                        className='w-8 px-1 ml-2 text-xs leading-tight text-gray-700 border border-purple-400 rounded focus:outline-none focus:bg-white'
                                        placeholder={t('homework:entergrade','Enter grade')}
                                        onChange={(e) => props.onChange(+e.target.value)}
                                        value={props.value}
                                      />
                                    )}
                                  />
                                  <span className="ml-8 text-xs font-bold tracking-wide text-gray-700">Comments: </span>
                                  <input
                                    name={`Responses[${index}].Comments`}
                                    ref={register}
                                    className='w-32 px-1 ml-2 text-xs leading-tight text-gray-700 border border-purple-400 rounded focus:outline-none focus:bg-white'
                                    type='text'
                                    defaultValue={responses[index].Comments}
                                    placeholder={t('homework:entercomments','Enter comments')}
                                  />
                                </div>
                              </div>
                            </div>
                            <Controller
                              control={control}
                              name={`Responses[${index}].ID`}
                              defaultValue={findResponseDefaultValue(question.ID)}
                              render={(props) => (
                                <input
                                  type='hidden'
                                  name={`Responses[${index}].ID`}
                                  className='outlie-none'
                                  onChange={(e) => props.onChange(+e.target.value)}
                                  value={props.value}
                                />
                              )}
                            />
                          </div>
                        </li>
                      );
                    case 2:
                      return (
                        <li key={question.ID} className="mt-5">
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
                                {question.Choices.map((choice, idx) => (
                                  <div key={idx} className='flex flex-row items-center px-4'>
                                    <div className='w-12 mx-2 outline-none'>
                                      {renderMultipleCorrect(questions[index], responses[index], idx)}
                                    </div>
                                    <div className='px-2 py-1 text-xs leading-tight text-gray-700 focus:outline-none focus:bg-white'>
                                      {choice}
                                    </div>
                                  </div>
                                ))}
                              </div>
                            </div>
                            <div className='flex flex-row items-center px-4 mt-2'>
                              <span className="ml-2 text-xs font-bold tracking-wide text-gray-700">Grade: </span>
                              <Controller
                                control={control}
                                name={`Responses[${index}].Grade`}
                                defaultValue={responses[index].Grade}
                                render={(props) => (
                                  <input
                                    type='text'
                                    name={`Responses[${index}].Grade`}
                                    className='w-8 px-1 ml-2 text-xs leading-tight text-gray-700 border border-purple-400 rounded focus:outline-none focus:bg-white'
                                    placeholder={t('homework:entergrade','Enter grade')}
                                    onChange={(e) => props.onChange(+e.target.value)}
                                    value={props.value}
                                  />
                                )}
                              />
                              <span className="ml-8 text-xs font-bold tracking-wide text-gray-700">Comments: </span>
                              <input
                                name={`Responses[${index}].Comments`}
                                ref={register}
                                className='w-32 px-1 ml-2 text-xs leading-tight text-gray-700 border border-purple-400 rounded focus:outline-none focus:bg-white'
                                type='text'
                                defaultValue={responses[index].Comments}
                                placeholder={t('homework:entercomments','Enter comments')}
                              />
                            </div>
                            <Controller
                              control={control}
                              name={`Responses[${index}].ID`}
                              defaultValue={findResponseDefaultValue(question.ID)}
                              render={(props) => (
                                <input
                                  type='hidden'
                                  name={`Responses[${index}].ID`}
                                  className='outlie-none'
                                  onChange={(e) => props.onChange(+e.target.value)}
                                  value={props.value}
                                />
                              )}
                            />
                          </div>
                        </li>
                      );
                    case 3:
                      return (
                        <li key={question.ID} className="mt-5">
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
                                <div className='w-12 mx-2 outline-none'>
                                  {responses[index].QuestionID === question.ID ?
                                      responses[index].TrueFalseResponse === questions[index].TrueFalseAnswer ?
                                      <span>&#9989;</span> : <span>&#10060;</span>
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

                                <div>
                                  <span className="ml-10 text-xs font-bold tracking-wide text-gray-700">Grade: </span>
                                  <Controller
                                    control={control}
                                    name={`Responses[${index}].Grade`}
                                    defaultValue={responses[index].Grade}
                                    render={(props) => (
                                      <input
                                        type='text'
                                        name={`Responses[${index}].Grade`}
                                        className='w-8 px-1 ml-2 text-xs leading-tight text-gray-700 border border-purple-400 rounded focus:outline-none focus:bg-white'
                                        placeholder={t('homework:entergrade','Enter grade')}
                                        onChange={(e) => props.onChange(+e.target.value)}
                                        value={props.value}
                                      />
                                    )}
                                  />
                                  <span className="ml-8 text-xs font-bold tracking-wide text-gray-700">Comments: </span>
                                  <input
                                    name={`Responses[${index}].Comments`}
                                    ref={register}
                                    className='w-32 px-1 mb-3 mr-4 text-xs leading-tight text-gray-700 border border-purple-400 rounded focus:outline-none focus:bg-white'
                                    type='text'
                                    defaultValue={responses[index].Comments}
                                    placeholder={t('homework:entercomments','Enter comments')}
                                  />
                                  <Controller
                                    control={control}
                                    name={`Responses[${index}].ID`}
                                    defaultValue={findResponseDefaultValue(question.ID)}
                                    render={(props) => (
                                      <input
                                        type='hidden'
                                        name={`Responses[${index}].ID`}
                                        className='outlie-none'
                                        onChange={(e) => props.onChange(+e.target.value)}
                                        value={props.value}
                                      />
                                    )}
                                  />
                                </div>
                              </div>
                            </div>
                          </div>
                        </li>
                      );
                    case 4:
                      return (
                        <li key={question.ID} className="mt-5">
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
                            <div className="px-4">
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
                              <div className='flex flex-row mt-1'>
                                <label className='block px-4 pt-1 mb-2 text-xs font-bold tracking-wide text-gray-700'>
                                  Student Response
                                </label>
                                <div className='px-2 py-1 text-xs leading-tight text-gray-700 border border-purple-400 rounded focus:outline-none focus:bg-white'>
                                  {responses[index].OpenResponse}
                                </div>
                              </div>
                              <div className="px-4">
                                <span className="text-xs font-bold tracking-wide text-gray-700">Grade: </span>
                                <Controller
                                  control={control}
                                  name={`Responses[${index}].Grade`}
                                  defaultValue={responses[index].Grade}
                                  render={(props) => (
                                    <input
                                      type='text'
                                      name={`Responses[${index}].Grade`}
                                      className='w-8 px-1 ml-2 text-xs leading-tight text-gray-700 border border-purple-400 rounded focus:outline-none focus:bg-white'
                                      placeholder={t('homework:entergrade','Enter grade')}
                                      onChange={(e) => props.onChange(+e.target.value)}
                                      value={props.value}
                                    />
                                  )}
                                />
                                <span className="ml-8 text-xs font-bold tracking-wide text-gray-700">Comments: </span>
                                <input
                                  name={`Responses[${index}].Comments`}
                                  ref={register}
                                  className='w-32 px-1 mb-3 mr-4 text-xs leading-tight text-gray-700 border border-purple-400 rounded focus:outline-none focus:bg-white'
                                  type='text'
                                  defaultValue={responses[index].Comments}
                                  placeholder={t('homework:entercomments','Enter comments')}
                                />
                                <Controller
                                  control={control}
                                  name={`Responses[${index}].ID`}
                                  defaultValue={findResponseDefaultValue(question.ID)}
                                  render={(props) => (
                                    <input
                                      type='hidden'
                                      name={`Responses[${index}].ID`}
                                      className='outlie-none'
                                      onChange={(e) => props.onChange(+e.target.value)}
                                      value={props.value}
                                    />
                                  )}
                                />
                              </div>
                            </div>
                          </div>
                        </li>
                      );
                    default:
                      return null
                }})}
                </ul>

                <div className='flex justify-center mt-5'>
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