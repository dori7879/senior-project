import { useSelector } from 'react-redux'
import DateTimePicker from "react-datetime-picker";
import { Redirect } from 'react-router-dom'
import { useState, useEffect } from 'react'
import { useForm, Controller } from 'react-hook-form'
import { useTranslation } from 'react-i18next';
import UserService from '../services/user';
import AttendanceService from '../services/attendance';

const AttendanceForm = () => {
  const { t } = useTranslation(['translation', 'attendance']);
  const { register, handleSubmit, control } = useForm();
  const isLoggedIn = useSelector((state) => state.auth.isLoggedIn)
  
  const [successfull, setSuccessfull] = useState(false)
  const [ownedGroups, setOwnedGroups] = useState([])
  const [studentLink, setStudentLink] = useState("")
  const [teacherLink, setTeacherLink] = useState("")

  useEffect(()=> {
    if (isLoggedIn) {
      UserService.getProfile()
        .then(
          (response) => {
            if (response.data.OwnedGroups) {
              setOwnedGroups(response.data.OwnedGroups.Groups)
            }
          },
          (error) => {
            console.log(error.message)
          }
        )
    }
  }, [isLoggedIn])

  const onSubmit = (data) => {
    console.log(data);
    AttendanceService.createAttendance(data)
      .then(
        (response) => {
          setStudentLink(response.data.StudentLink)
          setTeacherLink(response.data.TeacherLink)
          setSuccessfull(true);
        },
        (error) => {
          alert(error.message);
        }
      )
  };

  if (successfull) {
    return <Redirect to={{
      pathname: "/link-attendance",
      state: { studentLink, teacherLink }
    }} />
  }

  return (
    <form
      onSubmit={handleSubmit(onSubmit)}
      className="w-3/4 p-4 bg-purple-300 border border-purple-300 rounded"
    >
      {isLoggedIn ? null : (
        <div className="flex flex-row items-center pb-2 items">
          <label className="block px-4 pt-1 mb-2 text-xs font-bold tracking-wide text-gray-700 uppercase">
          {t('attendance:fullname', 'Full name*')}
          </label>
          <input
            name="TeacherFullName"
            ref={register}
            className="px-2 py-1 text-xs leading-tight text-gray-700 border border-purple-400 rounded focus:outline-none focus:bg-white"
            type="text"
            placeholder="Enter your full name"
          />
        </div>
      )}
      <div className="flex flex-row items-center pb-2 items">
        <label className="block px-4 pt-1 mb-2 text-xs font-bold tracking-wide text-gray-700 uppercase">
        {t('attendance:submit', 'Set who can submit attendances*')}
        </label>
        <select
          name="Mode"
          ref={register}
          className="text-xs bg-purple-100 border border-purple-300"
        >
          <option value="all">{t('attendance:everyone', 'Everyone')}</option>
          <option value="registered">{t('attendance:registered', 'Registered Accounts')}</option>
        </select>
      </div>
      <div className="flex flex-row items-center pb-2 items">
        <label className="block px-4 pt-1 mb-2 text-xs font-bold tracking-wide text-gray-700 uppercase">
        {t('attendance:coursetitle', 'Course Title')}*
        </label>
        <input
          ref={register}
          name="CourseTitle"
          className="px-2 py-1 text-xs leading-tight text-gray-700 border border-purple-400 rounded focus:outline-none focus:bg-white"
          type="text"
          placeholder="Enter course title"
        />
      </div>
      {isLoggedIn ? (
        <div className="flex flex-row items-center pb-2 items">
          <label className="block px-4 pt-1 mb-2 text-xs font-bold tracking-wide text-gray-700 uppercase">
          {t('homework:group', 'Group')}
          </label>
          <Controller
            control={control}
            name={`GroupID`}
            defaultValue={0}
            render={(props) => (
              <select
                name="GroupID"
                className="px-2 py-1 text-xs leading-tight text-gray-700 border border-purple-400 rounded focus:outline-none focus:bg-white"
                onChange={(e) => props.onChange(+e.target.value)}
                value={props.value}
              >
                <option value={0} key={0}>Choose a group</option>
                {ownedGroups.map((group, idx) => {
                  return <option value={group.ID} key={group.ID}>{group.Title}</option>
                })}
              </select>
            )}
          />
        </div>
      ) : null}
      <div className="flex flex-row items-center pb-2 items">
        <label className="block px-4 pt-1 mb-2 text-xs font-bold tracking-wide text-gray-700 uppercase">
        {t('attendance:title', 'Title*')}
        </label>
        <input
          ref={register}
          name="Title"
          className="px-2 py-1 text-xs leading-tight text-gray-700 border border-purple-400 rounded focus:outline-none focus:bg-white"
          type="text"
          placeholder="Enter the attendance title"
        />
      </div>
      <div className="flex flex-row items-center pb-2 items">
        <label className="block px-4 pt-1 mb-2 text-xs font-bold tracking-wide text-gray-700 uppercase">
          {t('attendance:open', 'Open date and time*')}
        </label>
        <Controller
          as={DateTimePicker}
          control={control}
          value="selected"
          name="OpenedAt"
          onChange={(date) => date}
          defaultValue={new Date()}
        />
      </div>
      <div className="flex flex-row items-center pb-3 items">
        <label className="block px-4 pt-1 mb-2 text-xs font-bold tracking-wide text-gray-700 uppercase">
        {t('attendance:close', 'Close date and time*')} 
        </label>
        <Controller
          as={DateTimePicker}
          control={control}
          value="selected"
          name="ClosedAt"
          onChange={(date) => date}
          defaultValue={new Date()}
        />
      </div>
      <div className="flex justify-center">
        <button
          type="submit"
          className="relative flex justify-center w-full px-2 py-1 mb-2 text-sm font-medium leading-4 text-purple-200 transition duration-150 ease-in-out bg-purple-800 border border-transparent rounded-md hover:text-white hover:bg-purple-500 focus:outline-none"
        >
         {t('attendance:generate', 'Generate Link')}
        </button>
      </div>
    </form>
  )
}

export default AttendanceForm
