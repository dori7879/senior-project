import { Redirect, useParams } from "react-router-dom";
import { useSelector } from "react-redux";
import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { useForm } from "react-hook-form";
import Parser from "html-react-parser";
import CKEditor from "ckeditor4-react";
import Footer from "../components/Footer";
import Header from "../components/Header";
import HomeworkService from "../services/homework";

const SubmitHomework = () => {
  const { t } = useTranslation(["translation", "homework"]);
  const { randomStr } = useParams();
  const isLoggedIn = useSelector((state) => state.auth.isLoggedIn);
  const { register, handleSubmit, setValue } = useForm();

  // Homework related state
  const [courseTitle, setCourseTitle] = useState("");
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [mode, setMode] = useState("all");
  // eslint-disable-next-line no-unused-vars
  const [closeDate, setCloseDate] = useState(new Date());
  // eslint-disable-next-line no-unused-vars
  const [homeworkID, setHomeworkID] = useState(-1);

  // Form related state
  const [attachments, setAttachments] = useState([]);
  const [successful, setSuccessful] = useState(false);
  const [errorMsg, setErrorMsg] = useState("");

  const ckEditorRemoveTags = (data) => {
    const editedData = data.replace("<p>", "").replace("</p>", "");
    return editedData;
  };

  useEffect(() => {
    register("Response");

    HomeworkService.fetchStudentHomework(randomStr)
      .then((response) => {
        if (response.data) {
          setHomeworkID(response.data.ID);
          setMode(response.data.Mode);
          setCourseTitle(response.data.CourseTitle);
          setTitle(response.data.Title);
          setDescription(response.data.Content);
          setCloseDate(response.data.ClosedAt);
        }
      })
      .catch((error) => {
        setErrorMsg("You are logged in or not in the group");
        console.log(error.message);
      });
  }, [randomStr, register]);

  const onSubmit = (data) => {
    HomeworkService.submitHomework({
      ...data,
      HomeworkID: parseInt(homeworkID),
    }).then(
      (response) => {
        setSuccessful(true);
      },
      (error) => {
        alert(error.message);
      }
    );
  };

  if (mode === "registered" && !isLoggedIn) {
    return <Redirect to="/login" />;
  }

  const isEmptyDesc = description.trim() === "";
  const data = ckEditorRemoveTags(description);
  const isEmptyFiles = attachments.length === 0;

  return (
    <div>
      <Header />
      <div className="flex flex-col items-center min-h-screen px-4 py-2 bg-purple-100 justify-top sm:px-6 lg:px-8 ">
        <div className="flex flex-col items-center justify-center w-3/4 pb-4">
          <div className="pt-2 text-xl font-bold text-purple-900">
            {t("homework:homework", "Homework")}
          </div>
          <div className="flex flex-row items-center w-full border border-purple-300 rounded-t items ">
            <div className="w-full p-4 bg-purple-300 border border-purple-300 rounded-tl">
              <h2 className="block px-4 pt-1 mb-3 text-xs font-bold tracking-wide text-gray-700 uppercase">
                <strong>{t("homework:coursetitle", "Course Title")}:</strong>{" "}
                <span className="text-purple-900">{courseTitle}</span>
              </h2>
              <h2 className="block px-4 pt-1 mb-2 text-xs font-bold tracking-wide text-gray-700 uppercase">
                <strong>{t("homework:hwtitle", "Homework Title:")}</strong>{" "}
                <span className="text-purple-900">{title}</span>
              </h2>
            </div>
            <div>
              {isLoggedIn ? null : (
                <div className="flex flex-row items">
                  <label className="block px-4 pt-1 mb-4 text-xs font-bold tracking-wide text-gray-700 uppercase">
                    {t("homework:fullname", "Full Name*:")}
                  </label>
                  <input
                    name="StudentFullName"
                    ref={register}
                    className="px-1 mb-3 mr-4 text-xs leading-tight text-gray-700 border border-purple-400 rounded focus:outline-none focus:bg-white"
                    type="text"
                    placeholder={t(
                      "homework:entername",
                      "Enter your full name"
                    )}
                  />
                </div>
              )}
              <div className="flex flex-row items-center items">
                <h2 className="block px-4 pt-1 mb-2 text-xs font-bold tracking-wide text-gray-700 uppercase">
                  <strong>Closes At:</strong>{" "}
                  <span className="text-purple-900">
                    {closeDate.toString()}
                  </span>
                </h2>
              </div>
            </div>
          </div>
          <div className="w-full p-4 bg-purple-300 border border-purple-300 rounded-b">
            {isEmptyDesc ? null : (
              <div>
                <h2 className="block px-4 pt-1 mb-2 text-xs font-bold tracking-wide text-gray-700 uppercase">
                  <strong>{t("homework:description", "Description")}</strong>
                  <br></br>{" "}
                </h2>
                <p className="ml-4 text-purple-900">{Parser(data)}</p>
              </div>
            )}

            {errorMsg ? (
              <div className="flex items-center justify-center">
                <p className="text-red-600">{errorMsg}</p>
              </div>
            ) : null}

            {successful ? (
              <div className="form-group">
                <div className="w-full pt-2 text-xl font-bold text-center text-purple-900 border border-purple-300 rounded">
                  {t("homework:submitted", "Submitted")}
                </div>
              </div>
            ) : (
              <form onSubmit={handleSubmit(onSubmit)}>
                <div className="flex flex-col pb-2 mx-4">
                  <h1 className="block px-4 pt-1 mt-4 mb-2 text-xs font-bold tracking-wide text-center text-purple-900 uppercase">
                    {t("homework:answer", "answer")}
                  </h1>
                  <CKEditor
                    name="Response"
                    onChange={(e) => setValue("Response", e.editor.getData())}
                  />
                </div>
                {isEmptyFiles ? null : (
                  <h1>{t("homework:attachments", "Attachments")}</h1>
                )}
                <div className="flex flex-row items-center pb-2 ml-2 items">
                  <label className="block px-2 pt-1 mb-2 text-xs font-bold tracking-wide text-gray-700 uppercase">
                    {t("homework:files", "Attach files")}
                  </label>
                  <input
                    onChange={(e) => setAttachments(e.target.value)}
                    className="w-1/2 px-2 py-1 text-xs leading-tight text-gray-700 border border-purple-400 rounded focus:outline-none focus:bg-white"
                    type="file"
                    multiple
                  />
                </div>
                <div className="flex justify-center">
                  <button
                    type="submit"
                    className="relative flex justify-center w-full px-2 py-1 mb-2 text-sm font-medium leading-4 text-purple-200 transition duration-150 ease-in-out bg-purple-800 border border-transparent rounded-md hover:bg-purple-500 focus:outline-none"
                  >
                    {t("homework:submithw", "Submit")}
                  </button>
                </div>
              </form>
            )}
          </div>
        </div>
        <small className="">
          {t("homework:required", "* - Required field")}{" "}
        </small>
      </div>
      <Footer />
    </div>
  );
};

export default SubmitHomework;
