import { useDispatch, useSelector } from 'react-redux'

import CKEditor from "ckeditor4-react";
import DateTimePicker from "react-datetime-picker";
import { Redirect } from 'react-router-dom'
//import { createHomeworkPage } from '../actions/homework'
import { useState, useEffect } from 'react'
import {useForm, Controller} from 'react-hook-form'
import { useTranslation } from 'react-i18next';

const HwForm = () => {
  const { t } = useTranslation(['translation', 'homework']);
  const [successfull, setSuccessfull] = useState(false)
  const [isClicked, setIsClicked] = useState(false)
  const isLoggedIn = useSelector((state) => state.auth.isLoggedIn)
  const dispatch = useDispatch()
  const { register, handleSubmit, control, setValue } = useForm();


  const onSubmit = (data) => {
    console.log(data);
    alert(JSON.stringify(data));
  };

  useEffect(() => {
    register("Content");
  });
  /*const handleSubmit = (e) => {
    e.preventDefault()

    dispatch(
      createHomeworkPage(
        courseTitle,
        title,
        description,
        files,
        openDate,
        closeDate,
        fullName,
        mode
      )
    )
      .then(() => {
        setSuccessfull(true)
        setIsClicked(true)
      })
      .catch(() => {
        setSuccessfull(false)
        setIsClicked(false)
      })
  }*/

  if (isClicked) {
    return <Redirect to='/link' />
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
        {t('homework:coursetitle', 'Course Title*')}
        </label>
        <input
          ref={register}
          name="CourseTitle"
          className="px-2 py-1 text-xs leading-tight text-gray-700 border border-purple-400 rounded focus:outline-none focus:bg-white"
          type="text"
          placeholder="Enter course title"
        />
      </div>
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
          className="relative flex justify-center w-full px-2 py-1 mb-2 text-sm font-medium leading-4  hover:text-white text-purple-200 transition duration-150 ease-in-out bg-purple-800 border border-transparent rounded-md hover:bg-purple-500 focus:outline-none"
        >
         {t('homework:generate', 'Generate Link')}
        </button>
      </div>
    </form>
  )
}

export default HwForm
