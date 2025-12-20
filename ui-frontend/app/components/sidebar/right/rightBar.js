

import Widget from "./widget/widget";
import style from "./rightBar.module.css";
import FriendsSuggestions from "./friendsSuggestions/friendsSuggestions";

export default function RightBar() {
  return (
    <aside className={style.rightSideBar}>
      <Widget />
      <FriendsSuggestions />
    </aside>
  );
}
