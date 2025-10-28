package handlers

import (
	"fmt"
	"forum/internal/domain"
	"net/http"
	"strconv"
)

/*MARK: React
 */
func ReactHandler(w http.ResponseWriter, r *http.Request) {
	/*userID := r.Context().Value("userID").(int64)*/

	user, ok := r.Context().Value("user").(*domain.User)
	if !ok || user == nil {
		http.Error(w, "❌ unauthenticated user", http.StatusUnauthorized)
		return
	}

	targetType := r.URL.Query().Get("target_type")
	targetID, _ := strconv.ParseInt(r.URL.Query().Get("target_id"), 10, 64)
	value, _ := strconv.Atoi(r.URL.Query().Get("value"))

	if err := reactionService.React(user.ID, targetType, targetID, value); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	/*fmt.Fprintf(w, "reaction saved: user=%d target=%s/%d value=%d\n", user.ID, targetType, targetID, value)*/
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

/*MARK: RemoveReaction
 */
func RemoveReactionHandler(w http.ResponseWriter, r *http.Request) {
	/*userID := r.Context().Value("userID").(int64)*/

	user, ok := r.Context().Value("user").(*domain.User)
	if !ok || user == nil {
		http.Error(w, "❌ unauthenticated user", http.StatusUnauthorized)
		return
	}

	targetType := r.URL.Query().Get("target_type")
	targetID, _ := strconv.ParseInt(r.URL.Query().Get("target_id"), 10, 64)

	if err := reactionService.RemoveReaction(user.ID, targetType, targetID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	/*fmt.Fprintf(w, "reaction deleted: user=%d target=%s/%d\n", user.ID, targetType, targetID)*/
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

/*MARK: GetCounts
 */
func GetCountsHandler(w http.ResponseWriter, r *http.Request) {
	targetType := r.URL.Query().Get("target_type")
	targetID, _ := strconv.ParseInt(r.URL.Query().Get("target_id"), 10, 64)

	likes, dislikes, err := reactionService.GetReactionCounts(targetType, targetID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Target %s/%d: %d likes, %d dislikes\n", targetType, targetID, likes, dislikes)
}
