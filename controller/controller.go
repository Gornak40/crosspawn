package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const someCode = `#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <stdbool.h>

#define MAXS 5005000
#define S_MI 60
#define S_HO 3600
#define S_DAY 86400
#define MAXN 100100

int k, m;
char buf[MAXS];

// queue
int qul = MAXN, qur = MAXN;
int ququ[MAXN << 1];

#define QU_PUSH(x) ququ[qur++] = (x)
#define QU_SIZE() (qur - qul)
#define QU_TOP() ququ[qul]
#define QU_POP() ++qul

bool go(const char* s) {
	int _, mon, day, ho, mi, se;
	sscanf(s, "%d-%d-%d %d:%d:%d", &_, &mon, &day, &ho, &mi, &se);
	int tm = se + mi * S_MI + ho * S_HO + day * S_DAY;
	for (int jm = 1; jm < mon; ++jm) {
		switch (jm) {
		case 2:
			tm += 29 * S_DAY;
			break;
		case 4:
		case 6:
		case 9:
		case 11:
			tm += 30 * S_DAY;
			break;
		default:
			tm += 31 * S_DAY;
			break;
		}
	}
	while (QU_SIZE() && tm - QU_TOP() >= k) {
		QU_POP();
	}
	QU_PUSH(tm);
	return QU_SIZE() >= m;
}

int main() {
	scanf("%d %d\n", &k, &m);
	char *s = buf;
	size_t ss = sizeof(buf);
	char *ans = NULL;
	while (gets(s)) {
		if (ans == NULL && go(s)) {
			ans = malloc(strlen(s) + 1);
			strcpy(ans, s);
		}
	}
	if (ans == NULL) {
		printf("-1\n");
		return 0;
	}
	ans[19] = '\n';
	ans[20] = '\0';
	printf(ans);
}`

func Router(r *gin.Engine) {
	r.GET("/", index)
	r.GET("/codereview", codereview)
}

func index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Title": "Home",
	})
}

func codereview(c *gin.Context) {
	c.HTML(http.StatusOK, "codereview.html", gin.H{
		"Title":     "CodeReview",
		"CodeTitle": "amogus.cpp",
		"Code":      someCode,
	})
}
