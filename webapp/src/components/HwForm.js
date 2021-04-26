import { useSelector } from 'react-redux'

import CKEditor from "ckeditor4-react";
import DateTimePicker from "react-datetime-picker";
import { Redirect } from 'react-router-dom'
import { useState, useEffect } from 'react'
import { useForm, Controller } from 'react-hook-form'
import { useTranslation } from 'react-i18next';
import UserService from '../services/user';
import HomeworkService from '../services/homework';

const HwForm = () => {
  const { t } = useTranslation(['translation', 'homework']);
  const { register, handleSubmit, control, setValue } = useForm();
  const isLoggedIn = useSelector((state) => state.auth.isLoggedIn)
  
  const [successfull, setSuccessfull] = useState(false)
  const [ownedGroups, setOwnedGroups] = useState([])
  const [studentLink, setStudentLink] = useState("")
  const [teacherLink, setTeacherLink] = useState("")

  useEffect(()=> {
    register("Content");
    
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
  }, [register, isLoggedIn])

  const onSubmit = (data) => {
    console.log(data);
    // HomeworkService.createHomework(data)
    //   .then(
    //     (response) => {
    //       setStudentLink(response.data.StudentLink)
    //       setTeacherLink(response.data.TeacherLink)
    //       setSuccessfull(true);
    //     },
    //     (error) => {
    //       alert(error.message);
    //     }
    //   )
  };

  if (successfull) {
    return <Redirect to={{
      pathname: "/link",
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
          {t('homework:fullname', 'Full name*')}
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
        {t('homework:submit', 'Set who can submit homeworks*')}
        </label>
        <select
          name="Mode"
          ref={register}
          className="text-xs bg-purple-100 border border-purple-300"
        >
          <option value="all">{t('homework:everyone', 'Everyone')}</option>
          <option value="registered">{t('homework:registered', 'Registered Accounts')}</option>
        </select>
      </div>
      <div className="flex flex-row items-center pb-2 items">
        <label className="block px-4 pt-1 mb-2 text-xs font-bold tracking-wide text-gray-700 uppercase">
        {t('homework:coursetitle', 'Course Title')}*
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
        {t('homework:title', 'Title*')}
        </label>
        <input
          ref={register}
          name="Title"
          className="px-2 py-1 text-xs leading-tight text-gray-700 border border-purple-400 rounded focus:outline-none focus:bg-white"
          type="text"
          placeholder="Enter the homework title"
        />
      </div>
      <div className="flex flex-row pb-2">
        <label className="block px-4 pt-1 mb-2 text-xs font-bold tracking-wide text-gray-700 uppercase">
        {t('homework:description', 'Description')}
        </label>
        <CKEditor
          name="Content"
          onChange={(e) => setValue("Content", e.editor.getData())}
        />
      </div>
      {/*
      <div className='flex flex-row items-center pb-2 items'>
        <label className='block px-4 pt-1 mb-2 text-xs font-bold tracking-wide text-gray-700 uppercase'>
          Attach files
        </label>
        <input
          ref={register}
          className='w-1/2 px-2 py-1 text-xs leading-tight text-gray-700 border border-purple-400 rounded focus:outline-none focus:bg-white'
          name='Files'
          type='file'
          multiple
        />
      </div>
      */}
      <div className="flex flex-row items-center pb-2 items">
        <label className="block px-4 pt-1 mb-2 text-xs font-bold tracking-wide text-gray-700 uppercase">
          {t('homework:open', 'Open date and time*')}
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
        {t('homework:close', 'Close date and time*')} 
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
         {t('homework:generate', 'Generate Link')}
        </button>
      </div>
    </form>
  )
}

export default HwForm
